package v1

import (
	"context"

	"course-service/internal/domain"
	"course-service/internal/services/assignments"
	"course-service/internal/services/bankquestions"
	"course-service/internal/services/banks"
	"course-service/internal/services/courses"
	"course-service/internal/services/coursesectionitems"
	"course-service/internal/services/coursesections"
	"course-service/internal/services/lectures"

	"github.com/rs/zerolog"
)

type CoursesService interface {
	Create(ctx context.Context, params courses.CreateServiceParams) (domain.Course, error)
	Update(ctx context.Context, params courses.UpdateServiceParams) (domain.Course, error)
	Delete(ctx context.Context, id int64) error
	Get(ctx context.Context, id int64) (domain.Course, error)
	GetList(ctx context.Context, params courses.GetListServiceParams) ([]domain.Course, error)
	GetWithAllItems(ctx context.Context, params courses.GetWithAllItemsParams) (domain.CourseWithItems, error)
}

type CourseListenersService interface {
	SetListener(ctx context.Context, params courses.SetListenerServiceParams) error
	GetListenersList(ctx context.Context, params courses.GetListenersListServiceParams) ([]int64, error)
	DeleteListener(ctx context.Context, params courses.DeleteListenerServiceParams) error
}

type CourseSectionsService interface {
	Create(ctx context.Context, params coursesections.CreateServiceParams) (domain.CourseSection, error)
	Update(ctx context.Context, params coursesections.UpdateServiceParams) (domain.CourseSection, error)
	Delete(ctx context.Context, id int64) error
	Get(ctx context.Context, id int64) (domain.CourseSection, error)
	GetList(ctx context.Context, params coursesections.GetListServiceParams) ([]domain.CourseSection, error)
}

type CourseSectionItemsService interface {
	Create(ctx context.Context, params coursesectionitems.CreateServiceParams) (domain.CourseSectionItem, error)
	Update(ctx context.Context, params coursesectionitems.UpdateServiceParams) (domain.CourseSectionItem, error)
	Delete(ctx context.Context, id int64) error
	Get(ctx context.Context, id int64) (domain.CourseSectionItem, error)
	GetList(ctx context.Context, params coursesectionitems.GetListServiceParams) ([]domain.CourseSectionItem, error)
}

type LecturesService interface {
	Create(ctx context.Context, params lectures.CreateServiceParams) (domain.Lecture, error)
	Update(ctx context.Context, params lectures.UpdateServiceParams) (domain.Lecture, error)
	Delete(ctx context.Context, itemId int64) error
	Get(ctx context.Context, itemId int64) (domain.Lecture, error)
	GetList(ctx context.Context, params lectures.GetListServiceParams) ([]domain.Lecture, error)
}

type AssignmentsService interface {
	Create(ctx context.Context, params assignments.CreateServiceParams) (domain.Assignment, error)
	Update(ctx context.Context, params assignments.UpdateServiceParams) (domain.Assignment, error)
	Delete(ctx context.Context, itemId int64) error
	Get(ctx context.Context, itemId int64) (domain.Assignment, error)
	GetList(ctx context.Context, params assignments.GetListServiceParams) ([]domain.Assignment, error)
}

type BanksService interface {
	Create(ctx context.Context, params banks.CreateServiceParams) (domain.Bank, error)
	Update(ctx context.Context, params banks.UpdateServiceParams) (domain.Bank, error)
	Delete(ctx context.Context, id int64) error
	GetList(ctx context.Context, params banks.GetListServiceParams) ([]domain.Bank, error)
}

type BankQuestionsService interface {
	Create(ctx context.Context, params bankquestions.CreateServiceParams) (domain.BankQuestion, error)
	Update(ctx context.Context, params bankquestions.UpdateServiceParams) (domain.BankQuestion, error)
	BulkCreate(ctx context.Context, params bankquestions.BulkCreateServiceParams) error
	BulkUpdate(ctx context.Context, params bankquestions.BulkUpdateServiceParams) error
	GetList(ctx context.Context, params bankquestions.GetListServiceParams) ([]domain.BankQuestion, error)
}

type Config struct {
	CoursesService            CoursesService
	CourseListenersService    CourseListenersService
	CourseSectionsService     CourseSectionsService
	CourseSectionItemsService CourseSectionItemsService
	LecturesService           LecturesService
	AssignmentsService        AssignmentsService
	BanksService              BanksService
	BankQuestionsService      BankQuestionsService
	Logger                    *zerolog.Logger
}

type Handler struct {
	coursesService            CoursesService
	courseListenersService    CourseListenersService
	courseSectionsService     CourseSectionsService
	courseSectionItemsService CourseSectionItemsService
	lecturesService           LecturesService
	assignmentsService        AssignmentsService
	banksService              BanksService
	bankQuestionsService      BankQuestionsService

	logger *zerolog.Logger
}

func NewHandler(config Config) *Handler {
	return &Handler{
		coursesService:            config.CoursesService,
		courseListenersService:    config.CourseListenersService,
		courseSectionsService:     config.CourseSectionsService,
		courseSectionItemsService: config.CourseSectionItemsService,
		lecturesService:           config.LecturesService,
		assignmentsService:        config.AssignmentsService,
		banksService:              config.BanksService,
		bankQuestionsService:      config.BankQuestionsService,
		logger:                    config.Logger,
	}
}
