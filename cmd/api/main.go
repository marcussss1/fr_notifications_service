package main

import (
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	_ "github.com/lib/pq"
	"gitlab.com/fr5270937/notifications_service/internal/config"
	"gitlab.com/fr5270937/notifications_service/internal/middleware"
	"gitlab.com/fr5270937/notifications_service/internal/pkg/metrics"
	client_http "gitlab.com/fr5270937/notifications_service/internal/services/client/delivery/http"
	client_repository "gitlab.com/fr5270937/notifications_service/internal/services/client/repository"
	client_usecase "gitlab.com/fr5270937/notifications_service/internal/services/client/usecase"
	mailing_http "gitlab.com/fr5270937/notifications_service/internal/services/mailing/delivery/http"
	mailing_repository "gitlab.com/fr5270937/notifications_service/internal/services/mailing/repository"
	mailing_usecase "gitlab.com/fr5270937/notifications_service/internal/services/mailing/usecase"
	message_http "gitlab.com/fr5270937/notifications_service/internal/services/message/delivery/http"
	message_repository "gitlab.com/fr5270937/notifications_service/internal/services/message/repository"
	message_usecase "gitlab.com/fr5270937/notifications_service/internal/services/message/usecase"
	pending_mailing_usecase "gitlab.com/fr5270937/notifications_service/internal/services/pending_mailing/usecase"
	"gopkg.in/yaml.v3"
)

func init() {
	envPath := ".env"
	if err := godotenv.Load(envPath); err != nil {
		log.Fatal("No .env file found")
	}
}

func main() {
	yamlPath, exists := os.LookupEnv("YAML_PATH")
	if !exists {
		log.Fatal("Yaml path not found")
	}

	yamlFile, err := os.ReadFile(yamlPath)
	if err != nil {
		log.Fatal(err)
	}

	var config config.Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatal(err)
	}

	db, err := sqlx.Open(config.Postgres.DB, config.Postgres.ConnectionToDB)
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = db.Close() }()

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(10)

	metrics, err := metrics.NewMetricsService()
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()
	e.Use(middleware.LoggerMiddleware)
	e.Use(echoprometheus.NewMiddleware("app"))
	e.GET("/api/v1/metrics", echoprometheus.NewHandler())
	e.GET("/api/v1/docs", func(ctx echo.Context) error {
		return ctx.File("docs/index.html")
	})
	e.GET("/api/v1/swagger.yaml", func(ctx echo.Context) error {
		return ctx.File("docs/swagger.yaml")
	})

	clientRepository := client_repository.NewClientRepository(db)
	messageRepository := message_repository.NewMessageRepository(db)
	mailingRepository := mailing_repository.NewMailingRepository(db)

	clientUsecase := client_usecase.NewClientUsecase(clientRepository)
	messageUsecase := message_usecase.NewMessagesUsecase(metrics, clientRepository, messageRepository, mailingRepository)
	mailingUsecase := mailing_usecase.NewMailingUsecase(mailingRepository, messageUsecase)
	pendingMailingUsecase := pending_mailing_usecase.NewPendingMailingUsecase(messageUsecase, mailingRepository)

	client_http.NewClientHandler(e, clientUsecase)
	mailing_http.NewMailingHandler(e, mailingUsecase)
	message_http.NewMessagesHandler(e, messageUsecase)

	go pendingMailingUsecase.PendingMailingWatchdog()
	e.Logger.Fatal(e.Start(config.Server.Port))
}
