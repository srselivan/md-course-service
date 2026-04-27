package http

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
	v1 "course-service/internal/transport/http/v1"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/rs/zerolog"
)

type coursesService interface {
	Create(ctx context.Context, params courses.CreateServiceParams) (domain.Course, error)
	Update(ctx context.Context, params courses.UpdateServiceParams) (domain.Course, error)
	Delete(ctx context.Context, id int64) error
	Get(ctx context.Context, id int64) (domain.Course, error)
	GetList(ctx context.Context, params courses.GetListServiceParams) ([]domain.Course, error)
	GetWithAllItems(ctx context.Context, params courses.GetWithAllItemsParams) (domain.CourseWithItems, error)
}

type courseListenersService interface {
	SetListener(ctx context.Context, params courses.SetListenerServiceParams) error
	GetListenersList(ctx context.Context, params courses.GetListenersListServiceParams) ([]int64, error)
	DeleteListener(ctx context.Context, params courses.DeleteListenerServiceParams) error
}

type courseSectionsService interface {
	Create(ctx context.Context, params coursesections.CreateServiceParams) (domain.CourseSection, error)
	Update(ctx context.Context, params coursesections.UpdateServiceParams) (domain.CourseSection, error)
	Delete(ctx context.Context, id int64) error
	Get(ctx context.Context, id int64) (domain.CourseSection, error)
	GetList(ctx context.Context, params coursesections.GetListServiceParams) ([]domain.CourseSection, error)
}

type courseSectionItemsService interface {
	Create(ctx context.Context, params coursesectionitems.CreateServiceParams) (domain.CourseSectionItem, error)
	Update(ctx context.Context, params coursesectionitems.UpdateServiceParams) (domain.CourseSectionItem, error)
	Delete(ctx context.Context, id int64) error
	Get(ctx context.Context, id int64) (domain.CourseSectionItem, error)
	GetList(ctx context.Context, params coursesectionitems.GetListServiceParams) ([]domain.CourseSectionItem, error)
}

type lecturesService interface {
	Create(ctx context.Context, params lectures.CreateServiceParams) (domain.Lecture, error)
	Update(ctx context.Context, params lectures.UpdateServiceParams) (domain.Lecture, error)
	Delete(ctx context.Context, itemId int64) error
	Get(ctx context.Context, itemId int64) (domain.Lecture, error)
	GetList(ctx context.Context, params lectures.GetListServiceParams) ([]domain.Lecture, error)
}

type assignmentsService interface {
	Create(ctx context.Context, params assignments.CreateServiceParams) (domain.Assignment, error)
	Update(ctx context.Context, params assignments.UpdateServiceParams) (domain.Assignment, error)
	Delete(ctx context.Context, itemId int64) error
	Get(ctx context.Context, itemId int64) (domain.Assignment, error)
	GetList(ctx context.Context, params assignments.GetListServiceParams) ([]domain.Assignment, error)
}

type banksService interface {
	Create(ctx context.Context, params banks.CreateServiceParams) (domain.Bank, error)
	Update(ctx context.Context, params banks.UpdateServiceParams) (domain.Bank, error)
	Delete(ctx context.Context, id int64) error
	GetList(ctx context.Context, params banks.GetListServiceParams) ([]domain.Bank, error)
}

type bankQuestionsService interface {
	Create(ctx context.Context, params bankquestions.CreateServiceParams) (domain.BankQuestion, error)
	Update(ctx context.Context, params bankquestions.UpdateServiceParams) (domain.BankQuestion, error)
	BulkCreate(ctx context.Context, params bankquestions.BulkCreateServiceParams) error
	BulkUpdate(ctx context.Context, params bankquestions.BulkUpdateServiceParams) error
	GetList(ctx context.Context, params bankquestions.GetListServiceParams) ([]domain.BankQuestion, error)
}

type Config struct {
	Addr                      string
	CoursesService            coursesService
	CourseListenersService    courseListenersService
	CourseSectionsService     courseSectionsService
	CourseSectionItemsService courseSectionItemsService
	LecturesService           lecturesService
	AssignmentsService        assignmentsService
	BanksService              banksService
	BankQuestionsService      bankQuestionsService
	Logger                    *zerolog.Logger
}
type Server struct {
	app  *fiber.App
	addr string

	coursesService            coursesService
	courseListenersService    courseListenersService
	courseSectionsService     courseSectionsService
	courseSectionItemsService courseSectionItemsService
	lecturesService           lecturesService
	assignmentsService        assignmentsService
	banksService              banksService
	bankQuestionsService      bankQuestionsService

	logger *zerolog.Logger
}

func NewServer(config Config) *Server {
	s := &Server{
		app:                       fiber.New(),
		addr:                      config.Addr,
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

	s.init()

	return s
}

func (s *Server) Run() error {
	return s.app.Listen(
		s.addr,
		fiber.ListenConfig{DisableStartupMessage: true},
	)
}

func (s *Server) Shutdown() error {
	return s.app.Shutdown()
}

func (s *Server) init() {
	apiGroup := s.app.Group("/cs")
	apiGroup.Use(logger.New())

	v1Group := apiGroup.Group("/v1")
	v1Handler := v1.NewHandler(v1.Config{
		CoursesService:            s.coursesService,
		CourseListenersService:    s.courseListenersService,
		CourseSectionsService:     s.courseSectionsService,
		CourseSectionItemsService: s.courseSectionItemsService,
		LecturesService:           s.lecturesService,
		AssignmentsService:        s.assignmentsService,
		BanksService:              s.banksService,
		BankQuestionsService:      s.bankQuestionsService,
		Logger:                    s.logger,
	})

	{
		v1Handler.NewCoursesRoutes(v1Group)
		v1Handler.NewCourseSectionsRoutes(v1Group)
		v1Handler.NewCourseSectionItemsRoutes(v1Group)
		v1Handler.NewBanksRoutes(v1Group)
		v1Handler.NewBankQuestionsRoutes(v1Group)
	}
}
