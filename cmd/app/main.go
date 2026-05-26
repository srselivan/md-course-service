package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	assignmentservice "course-service/internal/services/assignments"
	bankquestionsservice "course-service/internal/services/bankquestions"
	banksservice "course-service/internal/services/banks"
	coursesservice "course-service/internal/services/courses"
	coursesectionitemsservice "course-service/internal/services/coursesectionitems"
	coursesectionsservice "course-service/internal/services/coursesections"
	testsservice "course-service/internal/services/tests"
	"course-service/internal/transport/http"
	"course-service/pkg/kafka"
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

	if err = postgres.RunMigrations(postgresConn.DB, cfg.Postgres.MigrationsPath, cfg.Postgres.DBName); err != nil {
		log.Fatal().Err(err).Msg("failed to run migrations")
	}
	log.Info().Msg("successfully ran migrations")

	var kafkaProducer *kafka.Producer
	if len(cfg.Kafka.Brokers) > 0 {
		kafkaProducer, err = kafka.NewProducer(kafka.ProducerConfig{
			Brokers: cfg.Kafka.Brokers,
			Logger:  log,
		})
		if err != nil {
			log.Fatal().Err(err).Msg("failed to create kafka producer")
		}
	}

	filesClient := coursesservice.NewFilesClient(kafkaProducer, log)
	coursesService := coursesservice.NewService(coursesservice.NewRepository(postgresConn), filesClient, log)
	courseSectionsService := coursesectionsservice.NewService(coursesectionsservice.NewRepository(postgresConn), log)
	courseSectionItemsService := coursesectionitemsservice.NewService(coursesectionitemsservice.NewRepository(postgresConn), log)
	banksService := banksservice.NewService(banksservice.NewRepository(postgresConn), log)
	bankQuestionsService := bankquestionsservice.NewService(bankquestionsservice.NewRepository(postgresConn), log)
	assignmentsService := assignmentservice.NewService(assignmentservice.NewRepository(postgresConn), log)
	testsService := testsservice.NewService(testsservice.NewRepository(postgresConn), log)

	pingers := []http.Pinger{postgres.NewHealth(postgresConn)}
	if kafkaProducer != nil {
		pingers = append(pingers, kafkaProducer)
	}

	appCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go testsService.RunExpiredAttemptsWorker(appCtx, time.Minute)

	httpServer := http.NewServer(http.Config{
		Addr:                      cfg.HTTPServer.Addr,
		Pingers:                   pingers,
		Logger:                    log,
		CoursesService:            coursesService,
		CourseListenersService:    coursesService,
		CourseSectionsService:     courseSectionsService,
		CourseSectionItemsService: courseSectionItemsService,
		BanksService:              banksService,
		BankQuestionsService:      bankQuestionsService,
		AssignmentsService:        assignmentsService,
		TestsService:              testsService,
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
	cancel()
	log.Info().Msg("shutting down course service...")

	if kafkaProducer != nil {
		kafkaProducer.Close()
		log.Info().Msg("kafka producer closed")
	}

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
