package tests

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

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

func (r *Repository) Create(ctx context.Context, params CreateRepoParams) (domain.Test, error) {
	const query = `
		insert into tests (
			course_id, course_section_item_id, title, description, available_from, available_to,
			duration_seconds, max_attempts, max_score, questions_count, generation_settings, bank_ids
		)
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		returning id, course_id, course_section_item_id, title, description, available_from, available_to,
			duration_seconds, max_attempts, max_score, questions_count, generation_settings, created_at, bank_ids`

	var test testModel
	err := r.db.GetContext(ctx, &test, query,
		params.CourseID,
		params.CourseSectionItemID,
		params.Title,
		params.Description,
		params.AvailableFrom,
		params.AvailableTo,
		params.DurationSeconds,
		params.MaxAttempts,
		params.MaxScore,
		params.QuestionsCount,
		params.GenerationSettings,
		pq.Array(params.BankIDs),
	)
	if err != nil {
		return domain.Test{}, fmt.Errorf("insert test: %w", err)
	}
	return test.toDomain(), nil
}

func (r *Repository) Update(ctx context.Context, params UpdateRepoParams) (domain.Test, error) {
	const query = `
		update tests
		set course_id = $2,
			course_section_item_id = $3,
			title = $4,
			description = $5,
			available_from = $6,
			available_to = $7,
			duration_seconds = $8,
			max_attempts = $9,
			max_score = $10,
			questions_count = $11,
			generation_settings = $12,
			bank_ids = $13
		where id = $1
		returning id, course_id, course_section_item_id, title, description, available_from, available_to,
			duration_seconds, max_attempts, max_score, questions_count, generation_settings, created_at, bank_ids`

	var test testModel
	err := r.db.GetContext(ctx, &test, query,
		params.ID,
		params.CourseID,
		params.CourseSectionItemID,
		params.Title,
		params.Description,
		params.AvailableFrom,
		params.AvailableTo,
		params.DurationSeconds,
		params.MaxAttempts,
		params.MaxScore,
		params.QuestionsCount,
		params.GenerationSettings,
		pq.Array(params.BankIDs),
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Test{}, repository.ErrNotFound
		}
		return domain.Test{}, fmt.Errorf("update test: %w", err)
	}
	return test.toDomain(), nil
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	const query = `delete from tests where id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete test: %w", err)
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

func (r *Repository) Get(ctx context.Context, id int64) (domain.Test, error) {
	const query = `
		select id, course_id, course_section_item_id, title, description, available_from, available_to,
			duration_seconds, max_attempts, max_score, questions_count, generation_settings, created_at, bank_ids
		from tests
		where id = $1`

	var test testModel
	err := r.db.GetContext(ctx, &test, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Test{}, repository.ErrNotFound
		}
		return domain.Test{}, fmt.Errorf("get test: %w", err)
	}
	return test.toDomain(), nil
}

func (r *Repository) GetList(ctx context.Context, courseID *int64) ([]domain.Test, error) {
	const (
		listAllQuery = `
			select id, course_id, course_section_item_id, title, description, available_from, available_to,
				duration_seconds, max_attempts, max_score, questions_count, generation_settings, created_at, bank_ids
			from tests
			order by created_at desc`
		listByCourseQuery = `
			select id, course_id, course_section_item_id, title, description, available_from, available_to,
				duration_seconds, max_attempts, max_score, questions_count, generation_settings, created_at, bank_ids
			from tests
			where course_id = $1
			order by created_at desc`
	)

	var list []testModel
	var err error
	if courseID != nil {
		err = r.db.SelectContext(ctx, &list, listByCourseQuery, *courseID)
	} else {
		err = r.db.SelectContext(ctx, &list, listAllQuery)
	}
	if err != nil {
		return nil, fmt.Errorf("get tests list: %w", err)
	}

	result := make([]domain.Test, 0, len(list))
	for _, test := range list {
		result = append(result, test.toDomain())
	}
	return result, nil
}

func (r *Repository) GetAttempts(ctx context.Context, params GetAttemptsRepoParams) ([]domain.TestAttempt, error) {
	const (
		getAttemptsBaseQuery = `
			select id, test_id, user_id, attempt_number, status, grade_status, started_at, expires_at, completed_at, total_score
			from test_attempts`
		getAttemptsOrderBy = ` order by started_at desc`
	)

	args := []any{}
	filters := make([]string, 0, 3)
	if params.TestID != nil {
		args = append(args, *params.TestID)
		filters = append(filters, fmt.Sprintf("test_id = $%d", len(args)))
	}
	if params.UserID != nil {
		args = append(args, *params.UserID)
		filters = append(filters, fmt.Sprintf("user_id = $%d", len(args)))
	}
	if params.Status != nil {
		args = append(args, *params.Status)
		filters = append(filters, fmt.Sprintf("status = $%d", len(args)))
	}

	query := getAttemptsBaseQuery
	if len(filters) > 0 {
		query += " where " + strings.Join(filters, " and ")
	}
	query += getAttemptsOrderBy

	var attempts []testAttemptModel
	if err := r.db.SelectContext(ctx, &attempts, query, args...); err != nil {
		return nil, fmt.Errorf("get attempts: %w", err)
	}
	return testAttemptModels(attempts).toDomain(), nil
}

func (r *Repository) CountAttempts(ctx context.Context, testID, userID int64) (int, error) {
	const query = `
		select count(*)
		from test_attempts
		where test_id = $1 and user_id = $2`

	var count int
	if err := r.db.GetContext(ctx, &count, query, testID, userID); err != nil {
		return 0, fmt.Errorf("count attempts: %w", err)
	}
	return count, nil
}

func (r *Repository) SelectBankQuestionIDs(ctx context.Context, params SelectBankQuestionIDsRepoParams) ([]int64, error) {
	const query = `
		select bq.id
		from bank_question bq
		join tests t on bq.bank_id = any(t.bank_ids)
		where t.id = $1
		order by bq.id`

	args := []any{params.TestID}
	fullQuery := query
	if params.Limit > 0 {
		args = append(args, params.Limit)
		fullQuery += fmt.Sprintf(" limit $%d", len(args))
	}

	var ids []int64
	if err := r.db.SelectContext(ctx, &ids, fullQuery, args...); err != nil {
		return nil, fmt.Errorf("select bank question ids: %w", err)
	}
	return ids, nil
}

func (r *Repository) GetQuestionsWithAnswers(ctx context.Context, questionIDs []int64) ([]domain.Question, error) {
	const query = `
		select id, bank_id, type, text, points
		from bank_question
		where id = any($1)`

	if len(questionIDs) == 0 {
		return nil, nil
	}

	var rows []bankQuestionRow
	if err := r.db.SelectContext(ctx, &rows, query, pq.Array(questionIDs)); err != nil {
		return nil, fmt.Errorf("get bank questions: %w", err)
	}

	answersByQuestion, err := r.loadBankAnswers(ctx, questionIDs)
	if err != nil {
		return nil, err
	}

	result := make([]domain.Question, 0, len(rows))
	for _, row := range rows {
		question := row.toDomainQuestion()
		if question.Type != domain.TestingQuestionTypeText {
			question.Answers = answersByQuestion[question.ID]
		}
		result = append(result, question)
	}
	return result, nil
}

func (r *Repository) loadBankAnswers(ctx context.Context, questionIDs []int64) (map[int64][]domain.QuestionAnswer, error) {
	const query = `
		select id, question_id, text, is_correct
		from bank_answers
		where question_id = any($1)
		order by question_id, id`

	var answers []questionAnswerModel
	if err := r.db.SelectContext(ctx, &answers, query, pq.Array(questionIDs)); err != nil {
		return nil, fmt.Errorf("get bank answers: %w", err)
	}

	result := make(map[int64][]domain.QuestionAnswer, len(questionIDs))
	for _, a := range answers {
		result[a.QuestionID] = append(result[a.QuestionID], domain.QuestionAnswer{
			ID:         a.ID,
			QuestionID: a.QuestionID,
			Text:       a.Text,
			IsCorrect:  a.IsCorrect,
		})
	}
	return result, nil
}

func (r *Repository) CreateAttempt(ctx context.Context, params CreateAttemptRepoParams) (domain.TestAttempt, error) {
	var attempt testAttemptModel
	err := r.withTx(ctx, func(tx *sqlx.Tx) error {
		const query = `
			insert into test_attempts (test_id, user_id, attempt_number, status, grade_status, started_at, expires_at)
			values ($1, $2, $3, $4, $5, $6, $7)
			returning id, test_id, user_id, attempt_number, status, grade_status, started_at, expires_at, completed_at, total_score`

		if err := tx.GetContext(ctx, &attempt, query,
			params.TestID,
			params.UserID,
			params.AttemptNumber,
			domain.TestAttemptStatusInProgress,
			domain.TestGradeStatusNone,
			params.StartedAt,
			params.ExpiresAt,
		); err != nil {
			return fmt.Errorf("insert attempt: %w", err)
		}

		if len(params.Questions) == 0 {
			return nil
		}

		attemptQuestionIDs, err := bulkInsertAttemptQuestions(ctx, tx, attempt.ID, params.Questions)
		if err != nil {
			return err
		}

		return bulkInsertAttemptQuestionAnswers(ctx, tx, attemptQuestionIDs, params.Questions)
	})
	if err != nil {
		return domain.TestAttempt{}, err
	}
	return attempt.toDomain(), nil
}

func (r *Repository) GetAttempt(ctx context.Context, attemptID int64) (domain.TestAttempt, error) {
	const query = `
		select id, test_id, user_id, attempt_number, status, grade_status, started_at, expires_at, completed_at, total_score
		from test_attempts
		where id = $1`

	var attempt testAttemptModel
	err := r.db.GetContext(ctx, &attempt, query, attemptID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.TestAttempt{}, repository.ErrNotFound
		}
		return domain.TestAttempt{}, fmt.Errorf("get attempt: %w", err)
	}
	return attempt.toDomain(), nil
}

func (r *Repository) GetAttemptQuestion(ctx context.Context, attemptQuestionID int64) (domain.AttemptQuestion, error) {
	const query = `
		select id, attempt_id, question_id, order_index, question_type, question_text, points, score_awarded, instructor_comment
		from attempt_questions
		where id = $1`

	var row attemptQuestionModel
	err := r.db.GetContext(ctx, &row, query, attemptQuestionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.AttemptQuestion{}, repository.ErrNotFound
		}
		return domain.AttemptQuestion{}, fmt.Errorf("get attempt question: %w", err)
	}
	return row.toDomain(), nil
}

func (r *Repository) GetAttemptState(ctx context.Context, params GetAttemptStateRepoParams) (domain.AttemptState, error) {
	const (
		countQuestionsQuery = `
			select count(*) from attempt_questions where attempt_id = $1`
		listQuestionsQuery = `
			select id, attempt_id, question_id, order_index, question_type, question_text, points, score_awarded, instructor_comment
			from attempt_questions
			where attempt_id = $1
			order by order_index
			limit $2 offset $3`
	)

	attempt, err := r.GetAttempt(ctx, params.AttemptID)
	if err != nil {
		return domain.AttemptState{}, err
	}

	var total int
	if err = r.db.GetContext(ctx, &total, countQuestionsQuery, params.AttemptID); err != nil {
		return domain.AttemptState{}, fmt.Errorf("count attempt questions: %w", err)
	}

	var rows []attemptQuestionModel
	if err = r.db.SelectContext(ctx, &rows, listQuestionsQuery, params.AttemptID, params.Limit, params.Offset); err != nil {
		return domain.AttemptState{}, fmt.Errorf("get attempt questions: %w", err)
	}

	questions, err := r.buildAttemptQuestionStates(ctx, rows, params.IncludeCorrect)
	if err != nil {
		return domain.AttemptState{}, err
	}

	return domain.AttemptState{
		Attempt:   attempt,
		Questions: questions,
		Pagination: domain.AttemptPagination{
			Limit:  params.Limit,
			Offset: params.Offset,
			Total:  total,
		},
	}, nil
}

func (r *Repository) GetAttemptQuestionsForScoring(ctx context.Context, attemptID int64) ([]domain.AttemptQuestionState, error) {
	const query = `
		select id, attempt_id, question_id, order_index, question_type, question_text, points, score_awarded, instructor_comment
		from attempt_questions
		where attempt_id = $1
		order by order_index`

	var rows []attemptQuestionModel
	if err := r.db.SelectContext(ctx, &rows, query, attemptID); err != nil {
		return nil, fmt.Errorf("get attempt questions for scoring: %w", err)
	}
	return r.buildAttemptQuestionStates(ctx, rows, true)
}

func (r *Repository) GetAttemptScores(ctx context.Context, attemptID int64) ([]AttemptQuestionScore, error) {
	const query = `
		select id, score_awarded, instructor_comment
		from attempt_questions
		where attempt_id = $1`

	var rows []attemptScoreModel
	if err := r.db.SelectContext(ctx, &rows, query, attemptID); err != nil {
		return nil, fmt.Errorf("get attempt scores: %w", err)
	}
	result := make([]AttemptQuestionScore, 0, len(rows))
	for _, row := range rows {
		result = append(result, AttemptQuestionScore{
			AttemptQuestionID: row.ID,
			ScoreAwarded:      row.ScoreAwarded,
			InstructorComment: row.InstructorComment,
		})
	}
	return result, nil
}

func (r *Repository) SaveTextAnswer(ctx context.Context, params SaveTextAnswerRepoParams) error {
	return r.withTx(ctx, func(tx *sqlx.Tx) error {
		const (
			deleteAnswerQuery = `
				delete from attempt_answers where attempt_question_id = $1`
			insertAnswerQuery = `
				insert into attempt_answers (attempt_question_id, text_response)
				values ($1, $2)`
		)

		if _, err := tx.ExecContext(ctx, deleteAnswerQuery, params.AttemptQuestionID); err != nil {
			return fmt.Errorf("delete previous answer: %w", err)
		}
		if _, err := tx.ExecContext(ctx, insertAnswerQuery, params.AttemptQuestionID, params.TextResponse); err != nil {
			return fmt.Errorf("insert text answer: %w", err)
		}
		return nil
	})
}

func (r *Repository) SaveChoiceAnswers(ctx context.Context, params SaveChoiceAnswersRepoParams) error {
	return r.withTx(ctx, func(tx *sqlx.Tx) error {
		const (
			deleteAnswersQuery = `
				delete from attempt_answers where attempt_question_id = $1`
			insertAnswersQuery = `
				with inserted as (
					insert into attempt_answers (attempt_question_id, selected_answer_id)
					select $1, aqa.id
					from attempt_question_answers aqa
					where aqa.attempt_question_id = $1 and aqa.id = any($2)
					returning 1
				)
				select count(*) from inserted`
		)

		if _, err := tx.ExecContext(ctx, deleteAnswersQuery, params.AttemptQuestionID); err != nil {
			return fmt.Errorf("delete previous answers: %w", err)
		}
		if len(params.SelectedAnswerIDs) == 0 {
			return nil
		}

		var inserted int
		if err := tx.GetContext(ctx, &inserted, insertAnswersQuery,
			params.AttemptQuestionID, pq.Array(params.SelectedAnswerIDs)); err != nil {
			return fmt.Errorf("insert selected answers: %w", err)
		}
		if inserted != len(params.SelectedAnswerIDs) {
			return repository.ErrNotFound
		}
		return nil
	})
}

func (r *Repository) CompleteAttempt(ctx context.Context, params CompleteAttemptRepoParams) (domain.TestAttempt, error) {
	var attempt testAttemptModel
	err := r.withTx(ctx, func(tx *sqlx.Tx) error {
		const query = `
			update test_attempts
			set status = $2, grade_status = $3, completed_at = $4, total_score = $5
			where id = $1
			returning id, test_id, user_id, attempt_number, status, grade_status, started_at, expires_at, completed_at, total_score`

		if err := bulkUpdateAttemptQuestionScores(ctx, tx, params.AttemptID, params.Scores); err != nil {
			return err
		}
		if err := tx.GetContext(ctx, &attempt, query,
			params.AttemptID, params.Status, params.GradeStatus, params.CompletedAt, params.TotalScore); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return repository.ErrNotFound
			}
			return fmt.Errorf("complete attempt: %w", err)
		}
		return nil
	})
	if err != nil {
		return domain.TestAttempt{}, err
	}
	return attempt.toDomain(), nil
}

func (r *Repository) GradeAttempt(ctx context.Context, params GradeAttemptRepoParams) (domain.TestAttempt, error) {
	var attempt testAttemptModel
	err := r.withTx(ctx, func(tx *sqlx.Tx) error {
		const query = `
			update test_attempts
			set status = $2, grade_status = $3, completed_at = $4, total_score = $5
			where id = $1
			returning id, test_id, user_id, attempt_number, status, grade_status, started_at, expires_at, completed_at, total_score`

		if err := bulkUpdateAttemptQuestionScores(ctx, tx, params.AttemptID, params.Scores); err != nil {
			return err
		}
		if err := tx.GetContext(ctx, &attempt, query,
			params.AttemptID, domain.TestAttemptStatusCompleted, params.GradeStatus, params.CompletedAt, params.TotalScore); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return repository.ErrNotFound
			}
			return fmt.Errorf("grade attempt: %w", err)
		}
		return nil
	})
	if err != nil {
		return domain.TestAttempt{}, err
	}
	return attempt.toDomain(), nil
}

func (r *Repository) GetExpiredAttemptIDs(ctx context.Context, now time.Time) ([]int64, error) {
	const query = `
		select id
		from test_attempts
		where status = $1 and expires_at <= $2
		order by expires_at`

	var ids []int64
	if err := r.db.SelectContext(ctx, &ids, query, domain.TestAttemptStatusInProgress, now); err != nil {
		return nil, fmt.Errorf("get expired attempt ids: %w", err)
	}
	return ids, nil
}

func (r *Repository) buildAttemptQuestionStates(
	ctx context.Context,
	rows []attemptQuestionModel,
	includeCorrect bool,
) ([]domain.AttemptQuestionState, error) {
	if len(rows) == 0 {
		return nil, nil
	}

	attemptQuestionIDs := make([]int64, 0, len(rows))
	for _, row := range rows {
		attemptQuestionIDs = append(attemptQuestionIDs, row.ID)
	}

	answersByQuestion, err := r.loadAttemptQuestionAnswers(ctx, attemptQuestionIDs)
	if err != nil {
		return nil, err
	}
	responsesByQuestion, err := r.loadAttemptAnswers(ctx, attemptQuestionIDs)
	if err != nil {
		return nil, err
	}

	states := make([]domain.AttemptQuestionState, 0, len(rows))
	for _, row := range rows {
		attemptQuestion := row.toDomain()
		state := domain.AttemptQuestionState{
			AttemptQuestion: attemptQuestion,
			Question: domain.Question{
				ID:     row.QuestionID,
				Type:   row.QuestionType,
				Text:   row.QuestionText,
				Points: row.Points,
			},
			Response: responsesByQuestion[row.ID],
		}

		if row.QuestionType != domain.TestingQuestionTypeText {
			answers := answersByQuestion[row.ID]
			if includeCorrect {
				state.Question.Answers = answers
			} else {
				public := make([]domain.PublicQuestionAnswer, 0, len(answers))
				for _, a := range answers {
					public = append(public, domain.PublicQuestionAnswer{
						ID:         a.ID,
						QuestionID: a.QuestionID,
						Text:       a.Text,
					})
				}
				state.Answers = public
			}
		}

		states = append(states, state)
	}
	return states, nil
}

func (r *Repository) loadAttemptQuestionAnswers(
	ctx context.Context,
	attemptQuestionIDs []int64,
) (map[int64][]domain.QuestionAnswer, error) {
	const query = `
		select id, attempt_question_id as question_id, answer_text as text, is_correct
		from attempt_question_answers
		where attempt_question_id = any($1)
		order by attempt_question_id, id`

	var rows []questionAnswerModel
	if err := r.db.SelectContext(ctx, &rows, query, pq.Array(attemptQuestionIDs)); err != nil {
		return nil, fmt.Errorf("get attempt question answers: %w", err)
	}

	result := make(map[int64][]domain.QuestionAnswer, len(attemptQuestionIDs))
	for _, row := range rows {
		result[row.QuestionID] = append(result[row.QuestionID], domain.QuestionAnswer{
			ID:         row.ID,
			QuestionID: row.QuestionID,
			Text:       row.Text,
			IsCorrect:  row.IsCorrect,
		})
	}
	return result, nil
}

func (r *Repository) loadAttemptAnswers(
	ctx context.Context,
	attemptQuestionIDs []int64,
) (map[int64][]domain.AttemptAnswer, error) {
	const query = `
		select id, attempt_question_id, selected_answer_id, text_response
		from attempt_answers
		where attempt_question_id = any($1)
		order by attempt_question_id, id`

	var rows []attemptAnswerModel
	if err := r.db.SelectContext(ctx, &rows, query, pq.Array(attemptQuestionIDs)); err != nil {
		return nil, fmt.Errorf("get attempt answers: %w", err)
	}

	result := make(map[int64][]domain.AttemptAnswer, len(attemptQuestionIDs))
	for _, row := range rows {
		result[row.AttemptQuestionID] = append(result[row.AttemptQuestionID], domain.AttemptAnswer{
			ID:                row.ID,
			AttemptQuestionID: row.AttemptQuestionID,
			SelectedAnswerID:  row.SelectedAnswerID,
			TextResponse:      row.TextResponse,
		})
	}
	return result, nil
}

func (r *Repository) withTx(ctx context.Context, fn func(*sqlx.Tx) error) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	if err = fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}
	return nil
}

func bulkInsertAttemptQuestions(
	ctx context.Context,
	tx *sqlx.Tx,
	attemptID int64,
	questions []domain.Question,
) ([]int64, error) {
	const query = `
		insert into attempt_questions (
			attempt_id, question_id, order_index, question_type, question_text, points
		)
		select $1, qid, ord, qtype, qtext, pts
		from unnest($2::bigint[], $3::bigint[], $4::text[], $5::text[], $6::bigint[])
			as u(qid, ord, qtype, qtext, pts)
		order by ord
		returning id, order_index`

	if len(questions) == 0 {
		return nil, nil
	}

	questionIDs := make([]int64, len(questions))
	orderIndexes := make([]int64, len(questions))
	questionTypes := make([]string, len(questions))
	questionTexts := make([]string, len(questions))
	points := make([]int64, len(questions))
	for i, q := range questions {
		questionIDs[i] = q.ID
		orderIndexes[i] = int64(i + 1)
		questionTypes[i] = string(q.Type)
		questionTexts[i] = q.Text
		points[i] = int64(q.Points)
	}

	type insertedRow struct {
		ID         int64 `db:"id"`
		OrderIndex int   `db:"order_index"`
	}
	var inserted []insertedRow
	if err := tx.SelectContext(ctx, &inserted, query,
		attemptID,
		pq.Array(questionIDs),
		pq.Array(orderIndexes),
		pq.Array(questionTypes),
		pq.Array(questionTexts),
		pq.Array(points),
	); err != nil {
		return nil, fmt.Errorf("bulk insert attempt questions: %w", err)
	}
	if len(inserted) != len(questions) {
		return nil, fmt.Errorf("bulk insert attempt questions: expected %d rows, got %d", len(questions), len(inserted))
	}

	result := make([]int64, len(questions))
	for _, row := range inserted {
		idx := row.OrderIndex - 1
		if idx < 0 || idx >= len(result) {
			return nil, fmt.Errorf("bulk insert attempt questions: unexpected order_index %d", row.OrderIndex)
		}
		result[idx] = row.ID
	}
	return result, nil
}

func bulkInsertAttemptQuestionAnswers(
	ctx context.Context,
	tx *sqlx.Tx,
	attemptQuestionIDs []int64,
	questions []domain.Question,
) error {
	const query = `
		insert into attempt_question_answers (
			attempt_question_id, source_answer_id, answer_text, is_correct
		)
		select aq_id, src_id, txt, is_corr
		from unnest($1::bigint[], $2::bigint[], $3::text[], $4::bool[])
			as u(aq_id, src_id, txt, is_corr)`

	totalAnswers := 0
	for _, q := range questions {
		totalAnswers += len(q.Answers)
	}
	if totalAnswers == 0 {
		return nil
	}

	aqIDs := make([]int64, 0, totalAnswers)
	sourceIDs := make([]int64, 0, totalAnswers)
	texts := make([]string, 0, totalAnswers)
	correct := make([]bool, 0, totalAnswers)
	for i, q := range questions {
		for _, a := range q.Answers {
			aqIDs = append(aqIDs, attemptQuestionIDs[i])
			sourceIDs = append(sourceIDs, a.ID)
			texts = append(texts, a.Text)
			correct = append(correct, a.IsCorrect)
		}
	}

	if _, err := tx.ExecContext(ctx, query,
		pq.Array(aqIDs),
		pq.Array(sourceIDs),
		pq.Array(texts),
		pq.Array(correct),
	); err != nil {
		return fmt.Errorf("bulk insert attempt question answers: %w", err)
	}
	return nil
}

func bulkUpdateAttemptQuestionScores(
	ctx context.Context,
	tx *sqlx.Tx,
	attemptID int64,
	scores []AttemptQuestionScore,
) error {
	const query = `
		update attempt_questions aq
		set score_awarded = u.score,
			instructor_comment = u.comment
		from unnest($1::bigint[], $2::bigint[], $3::text[])
			as u(id, score, comment)
		where aq.id = u.id and aq.attempt_id = $4`

	if len(scores) == 0 {
		return nil
	}

	ids := make([]int64, len(scores))
	awarded := make([]int64, len(scores))
	comments := make([]sql.NullString, len(scores))
	for i, score := range scores {
		ids[i] = score.AttemptQuestionID
		awarded[i] = int64(score.ScoreAwarded)
		if score.InstructorComment != nil {
			comments[i] = sql.NullString{String: *score.InstructorComment, Valid: true}
		}
	}

	result, err := tx.ExecContext(ctx, query,
		pq.Array(ids),
		pq.Array(awarded),
		pq.Array(comments),
		attemptID,
	)
	if err != nil {
		return fmt.Errorf("bulk update attempt question scores: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("bulk update rows affected: %w", err)
	}
	if rows != int64(len(scores)) {
		return repository.ErrNotFound
	}
	return nil
}
