package bankquestions

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"course-service/internal/domain"
	"course-service/internal/repository"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, params CreateRepoParams) (domain.BankQuestion, error) {
	var question domain.BankQuestion
	err := repository.WithTx(ctx, r.db, func(tx *sqlx.Tx) error {
		q, err := insertBankQuestion(ctx, tx, params.BankId, params.Text, params.Type, params.Points)
		if err != nil {
			return err
		}
		answers, err := insertAnswers(ctx, tx, q.ID, params.Answers)
		if err != nil {
			return err
		}
		question = q.toDomain(answers)
		return nil
	})
	if err != nil {
		return domain.BankQuestion{}, err
	}
	return question, nil
}

func (r *Repository) BulkCreate(ctx context.Context, params BulkCreateRepoParams) error {
	return repository.WithTx(ctx, r.db, func(tx *sqlx.Tx) error {
		for _, q := range params.Questions {
			model, err := insertBankQuestion(ctx, tx, q.BankId, q.Text, q.Type, q.Points)
			if err != nil {
				return err
			}
			if _, err := insertAnswers(ctx, tx, model.ID, q.Answers); err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *Repository) Update(ctx context.Context, params UpdateRepoParams) (domain.BankQuestion, error) {
	var question domain.BankQuestion
	err := repository.WithTx(ctx, r.db, func(tx *sqlx.Tx) error {
		const (
			updateQuestionQuery = `
				update bank_question
				set bank_id = $2,
					text = $3,
					type = $4,
					points = $5
				where id = $1
				returning id, bank_id, text, type, points, created_at`
			deleteAnswersQuery = `delete from bank_answers where question_id = $1`
		)

		var q bankQuestionModel
		if err := tx.GetContext(ctx, &q, updateQuestionQuery,
			params.ID, params.BankId, params.Text, string(params.Type), params.Points); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return repository.ErrNotFound
			}
			return fmt.Errorf("update question: %w", err)
		}

		if _, err := tx.ExecContext(ctx, deleteAnswersQuery, params.ID); err != nil {
			return fmt.Errorf("delete answers: %w", err)
		}

		answers, err := insertAnswers(ctx, tx, q.ID, params.Answers)
		if err != nil {
			return err
		}
		question = q.toDomain(answers)
		return nil
	})
	if err != nil {
		return domain.BankQuestion{}, err
	}
	return question, nil
}

func (r *Repository) BulkUpdate(ctx context.Context, params BulkUpdateRepoParams) error {
	return repository.WithTx(ctx, r.db, func(tx *sqlx.Tx) error {
		const deleteQuestionsQuery = `delete from bank_question where bank_id = $1`

		if _, err := tx.ExecContext(ctx, deleteQuestionsQuery, params.Id); err != nil {
			return fmt.Errorf("delete questions: %w", err)
		}

		for _, q := range params.Questions {
			model, err := insertBankQuestion(ctx, tx, q.BankId, q.Text, q.Type, q.Points)
			if err != nil {
				return err
			}
			if _, err := insertAnswers(ctx, tx, model.ID, q.Answers); err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *Repository) GetList(ctx context.Context, params GetListRepoParams) (GetListRepoResult, error) {
	const (
		listQuestionsBaseQuery = `
			select id, bank_id, text, type, points, created_at
			from bank_question`
		countQuestionsBaseQuery = `select count(*) from bank_question`
		listAnswersQuery        = `
			select id, question_id, text, is_correct
			from bank_answers
			where question_id = any($1)
			order by id`
	)

	whereClause, args := buildListQuestionsFilter(params)
	countQuery := countQuestionsBaseQuery + whereClause
	listQuery := listQuestionsBaseQuery + whereClause + " order by id"

	argNum := len(args)
	listQuery += fmt.Sprintf(" limit $%d offset $%d", argNum+1, argNum+2)
	listArgs := append(append([]any{}, args...), params.Limit, params.Offset)

	var total int64
	if err := r.db.GetContext(ctx, &total, countQuery, args...); err != nil {
		return GetListRepoResult{}, fmt.Errorf("count questions: %w", err)
	}

	var questions []bankQuestionModel
	if err := r.db.SelectContext(ctx, &questions, listQuery, listArgs...); err != nil {
		return GetListRepoResult{}, fmt.Errorf("list questions: %w", err)
	}

	if len(questions) == 0 {
		return GetListRepoResult{
			Items: []domain.BankQuestion{},
			Total: total,
		}, nil
	}

	ids := make([]int64, 0, len(questions))
	for _, q := range questions {
		ids = append(ids, q.ID)
	}

	var answers []bankAnswerModel
	if err := r.db.SelectContext(ctx, &answers, listAnswersQuery, pq.Array(ids)); err != nil {
		return GetListRepoResult{}, fmt.Errorf("list answers: %w", err)
	}

	answersByQuestion := make(map[int64][]bankAnswerModel, len(questions))
	for _, a := range answers {
		answersByQuestion[a.QuestionId] = append(answersByQuestion[a.QuestionId], a)
	}

	result := make([]domain.BankQuestion, 0, len(questions))
	for _, q := range questions {
		result = append(result, q.toDomain(answersByQuestion[q.ID]))
	}

	return GetListRepoResult{
		Items: result,
		Total: total,
	}, nil
}

func buildListQuestionsFilter(params GetListRepoParams) (string, []any) {
	var (
		conds  []string
		args   []any
		argNum int
	)

	argNum++
	conds = append(conds, fmt.Sprintf("bank_id = $%d", argNum))
	args = append(args, params.BankId)

	if params.Type != nil && *params.Type != "" {
		argNum++
		conds = append(conds, fmt.Sprintf("type = $%d", argNum))
		args = append(args, *params.Type)
	}
	if params.Filter != nil && *params.Filter != "" {
		argNum++
		conds = append(conds, fmt.Sprintf("text ilike $%d", argNum))
		args = append(args, "%"+*params.Filter+"%")
	}

	return " where " + strings.Join(conds, " and "), args
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	const query = `delete from bank_question where id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete question: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete rows affected: %w", err)
	}
	if rows == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func insertBankQuestion(
	ctx context.Context,
	tx *sqlx.Tx,
	bankID int64,
	text string,
	questionType domain.TestingQuestionType,
	points int,
) (bankQuestionModel, error) {
	const query = `
		insert into bank_question (bank_id, text, type, points)
		values ($1, $2, $3, $4)
		returning id, bank_id, text, type, points, created_at`

	var model bankQuestionModel
	if err := tx.GetContext(ctx, &model, query, bankID, text, string(questionType), points); err != nil {
		return bankQuestionModel{}, fmt.Errorf("insert question: %w", err)
	}
	return model, nil
}

func insertAnswers(ctx context.Context, tx *sqlx.Tx, questionID int64, answers []BankAnswerDTO) ([]bankAnswerModel, error) {
	const query = `
		insert into bank_answers (question_id, text, is_correct)
		values ($1, $2, $3)
		returning id, question_id, text, is_correct`

	if len(answers) == 0 {
		return nil, nil
	}

	inserted := make([]bankAnswerModel, 0, len(answers))
	for _, a := range answers {
		var model bankAnswerModel
		if err := tx.GetContext(ctx, &model, query, questionID, a.Text, a.IsCorrect); err != nil {
			return nil, fmt.Errorf("insert answer: %w", err)
		}
		inserted = append(inserted, model)
	}
	return inserted, nil
}
