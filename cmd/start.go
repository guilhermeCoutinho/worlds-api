package cmd

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-pg/pg"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/guilhermeCoutinho/worlds-api/dal"
	"github.com/guilhermeCoutinho/worlds-api/handler"
	"github.com/guilhermeCoutinho/worlds-api/services"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type App struct {
	Router     *mux.Router
	AuthRouter *mux.Router
	DAL        *dal.DAL
	Services   *services.Services
	logger     *logrus.Logger
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start worlds api",
	Long:  `start worlds api`,
	Run: func(cmd *cobra.Command, args []string) {
		StartServer()
	},
}

func StartServer() {
	fmt.Println("Starting server...")
	worldsApp := NewApp()
	worldsApp.Run()
}

func init() {
	rootCmd.AddCommand(startCmd)
}

func NewApp() *App {
	config := viper.New()
	logger := logrus.New()
	db := initPg(logger)
	redisClient := initRedis(logger)
	dal := dal.NewDAL(db)

	eventPublisher := services.NewRedisEventPublisher(redisClient, logger)
	services := services.NewServices(config, dal, logger, eventPublisher)

	router := mux.NewRouter()
	authRouter := router.PathPrefix("/").Subrouter()

	app := &App{
		Router:     router,
		AuthRouter: authRouter,
		DAL:        dal,
		Services:   services,
		logger:     logger,
	}

	app.SetupMiddlewares()
	app.SetupRoutes()
	return app
}

func initPg(logger *logrus.Logger) *pg.DB {
	logger.WithField("method", "initPg").Info("Connecting to postgres")
	pgURL := os.Getenv("PG_URL")
	if pgURL == "" {
		logger.WithField("method", "initPg").Fatal("PG_URL environment variable is not set")
	}

	opts, err := pg.ParseURL(pgURL)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"method": "initPg",
			"error":  err,
		}).Fatal("Failed to parse PG_URL")
	}

	db := pg.Connect(opts)
	_, err = db.Exec("SELECT 1")
	if err != nil {
		logger.WithFields(logrus.Fields{
			"method": "initPg",
			"error":  err,
		}).Fatal("Failed to connect to Postgres")
	}
	logger.WithField("method", "initPg").Info("Connected to Postgres")
	return db
}

func initRedis(logger *logrus.Logger) *redis.Client {
	logger.WithField("method", "initRedis").Info("Connecting to Redis")
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		logger.WithField("method", "initRedis").Warn("REDIS_URL environment variable is not set, using default localhost:6379")
		redisURL = "redis://localhost:6379"
	}

	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"method": "initRedis",
			"error":  err,
		}).Fatal("Failed to parse REDIS_URL")
	}

	client := redis.NewClient(opt)

	_, err = client.Ping(client.Context()).Result()
	if err != nil {
		logger.WithFields(logrus.Fields{
			"method": "initRedis",
			"error":  err,
		}).Fatal("Failed to connect to Redis")
	}

	logger.WithField("method", "initRedis").Info("Connected to Redis")
	return client
}

func (a *App) SetupRoutes() {
	handlers := handler.NewHandlers(a.Services)
	handlers.RegisterRoutes(a.Router)
	handlers.RegisterAuthenticatedRoutes(a.AuthRouter)
}

func (a *App) SetupMiddlewares() {
	authMiddleware := handler.NewAuthMiddleware(a.logger)
	a.AuthRouter.Use(authMiddleware.Authenticate)

	a.Router.Use(mux.CORSMethodMiddleware(a.Router))
	a.AuthRouter.Use(mux.CORSMethodMiddleware(a.AuthRouter))
}

func (a *App) Run() {
	logger := a.logger.WithField("method", "Run")
	logger.Info("Starting server on port 8080")
	err := http.ListenAndServe(":8080", a.Router)
	if err != nil {
		logger.Fatal(err)
	}
}
