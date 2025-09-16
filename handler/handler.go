package handler

import (
	"context"
	"net/http"
	"reflect"

	"github.com/gorilla/mux"
	"github.com/guilhermeCoutinho/worlds-api/services"
	"github.com/guilhermeCoutinho/worlds-api/utils"
	"github.com/sirupsen/logrus"
)

type Handlers struct {
	WorldsHandler      *WorldsHandler
	HealthcheckHandler *HealthcheckHandler
}

func LoggedHandler(handler func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId := ""
		if r.Context().Value("userID") != nil {
			userId = r.Context().Value("userID").(string)
		}

		logger := logrus.WithFields(logrus.Fields{
			"handler": r.URL.Path,
			"method":  r.Method,
			"userID":  userId,
		})

		logger.Info("Handling request")

		ctx := context.WithValue(r.Context(), utils.LoggerCtxKey, logger)
		err := handler(w, r.WithContext(ctx))
		if err != nil {
			logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}

func NewHandlers(services *services.Services) *Handlers {
	worldsHandler := NewWorldsHandler(services)
	healthcheckHandler := NewHealthcheckHandler()
	return &Handlers{
		WorldsHandler:      worldsHandler,
		HealthcheckHandler: healthcheckHandler,
	}
}

func (h *Handlers) RegisterRoutes(r *mux.Router) {
	logger := logrus.WithField("method", "RegisterRoutes")
	val := reflect.ValueOf(h).Elem()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if field.IsNil() {
			continue
		}
		method := field.MethodByName("RegisterHandler")
		if method.IsValid() {
			method.Call([]reflect.Value{reflect.ValueOf(r)})
			logger.WithField("handler", field.Type().Name()).Info("Registered handler")
		}
	}
}

func (h *Handlers) RegisterAuthenticatedRoutes(r *mux.Router) {
	logger := logrus.WithField("method", "RegisterAuthenticatedRoutes")
	val := reflect.ValueOf(h).Elem()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if field.IsNil() {
			continue
		}
		method := field.MethodByName("RegisterAuthenticatedHandler")
		if method.IsValid() {
			method.Call([]reflect.Value{reflect.ValueOf(r)})
			logger.WithField("handler", field.Type().Name()).Info("Registered authenticated handler")
		}
	}
}
