package main

import (
	"os"
	"os/signal"
	"syscall"

	"course-service/internal/repository/bankquestions"
	"course-service/internal/repository/banks"
	"course-service/internal/repository/courses"
	"course-service/internal/repository/coursesectionitems"
	"course-service/internal/repository/coursesections"
	bankquestionsservice "course-service/internal/services/bankquestions"
	banksservice "course-service/internal/services/banks"
	coursesservice "course-service/internal/services/courses"
	coursesectionitemsservice "course-service/internal/services/coursesectionitems"
	coursesectionsservice "course-service/internal/services/coursesections"
	"course-service/internal/transport/http"
	"course-service/pkg/gorm"
	"course-service/pkg/logger"
	"course-service/pkg/postgres"

	"course-service/internal/config"
)

func main() {
	const appName = "course-service"

	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	log, err := logger.New(logger.Config{
		Level:       cfg.Logger.Level,
		AppName:     appName,
		LogFilePath: cfg.Logger.FilePath,
	})
	if err != nil {
		panic(err)
	}

	postgresConn, err := postgres.New(postgres.Config{
		Host:     cfg.Postgres.Host,
		Port:     cfg.Postgres.Port,
		User:     cfg.Postgres.User,
		Password: cfg.Postgres.Password,
		DBName:   cfg.Postgres.DBName,
		SSLMode:  cfg.Postgres.SSLMode,
		AppName:  appName,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to postgres")
	}

	gormDb, err := gorm.New(postgresConn)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create orm")
	}

	coursesRepo := courses.NewPostgresRepo(gormDb)
	courseSectionsRepo := coursesections.NewPostgresRepo(gormDb)
	courseSectionsItemsRepo := coursesectionitems.NewPostgresRepo(gormDb)
	banksRepo := banks.NewPostgresRepo(gormDb)
	bankQuestionsRepo := bankquestions.NewPostgresRepo(gormDb)

	coursesService := coursesservice.NewService(coursesRepo, log)
	courseSectionsService := coursesectionsservice.NewService(courseSectionsRepo, log)
	courseSectionItemsService := coursesectionitemsservice.NewService(courseSectionsItemsRepo, log)
	banksService := banksservice.NewService(banksRepo, log)
	bankQuestionsService := bankquestionsservice.NewService(bankQuestionsRepo, log)

	httpServer := http.NewServer(http.Config{
		Addr:                      cfg.HTTPServer.Addr,
		Logger:                    log,
		CoursesService:            coursesService,
		CourseListenersService:    coursesService,
		CourseSectionsService:     courseSectionsService,
		CourseSectionItemsService: courseSectionItemsService,
		BanksService:              banksService,
		BankQuestionsService:      bankQuestionsService,
	})

	go func() {
		if err = httpServer.Run(); err != nil {
			log.Error().Err(err).Msg("failed to start http server")
		}
	}()

	log.Info().Msgf("http server started at %s", cfg.HTTPServer.Addr)
	log.Info().Msg("course service started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGKILL)
	<-quit
	log.Info().Msg("shutting down course service...")

	if err = postgresConn.Close(); err != nil {
		log.Error().Err(err).Msg("failed to close postgres connection")
	}
	log.Info().Msg("postgres connection closed")

	if err = httpServer.Shutdown(); err != nil {
		log.Error().Err(err).Msg("failed to shutdown http server")
	}
	log.Info().Msg("http server shutdown complete")

	log.Info().Msg("course service shutdown complete")
}
