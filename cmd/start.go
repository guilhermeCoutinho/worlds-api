package cmd

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-pg/pg"
	"github.com/gorilla/mux"
	"github.com/guilhermeCoutinho/worlds-api/dal"
	"github.com/guilhermeCoutinho/worlds-api/handler"
	"github.com/guilhermeCoutinho/worlds-api/services"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

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

type App struct {
	Router     *mux.Router
	AuthRouter *mux.Router
	DAL        *dal.DAL
	Services   *services.Services
	logger     *logrus.Logger
}

func NewApp() *App {
	config := viper.New()
	logger := logrus.New()
	db := initPg(logger)
	dal := dal.NewDAL(db)
	services := services.NewServices(config, dal, logger)

	router := mux.NewRouter()
	authRouter := router.PathPrefix("/").Subrouter()
	authMiddleware := handler.NewAuthMiddleware(logger)
	authRouter.Use(authMiddleware.Authenticate)

	app := &App{
		Router:     router,
		AuthRouter: authRouter,
		DAL:        dal,
		Services:   services,
		logger:     logger,
	}
	app.SetupRoutes()
	app.SetupMiddlewares()
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

func (a *App) SetupRoutes() {
	handlers := handler.NewHandlers(a.Services)
	handlers.RegisterRoutes(a.Router)
	handlers.RegisterAuthenticatedRoutes(a.AuthRouter)
}

func (a *App) SetupMiddlewares() {
	a.Router.Use(mux.CORSMethodMiddleware(a.Router))
}

func (a *App) Run() {
	logger := a.logger.WithField("method", "Run")
	logger.Info("Starting server on port 8080")
	err := http.ListenAndServe(":8080", a.Router)
	if err != nil {
		logger.Fatal(err)
	}
}
