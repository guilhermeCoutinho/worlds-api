package handler

import (
	"net/http"
	"reflect"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	"github.com/guilhermeCoutinho/worlds-api/services"
	"github.com/guilhermeCoutinho/worlds-api/utils"
	"github.com/sirupsen/logrus"
)

type Handlers struct {
	WorldsHandler      *WorldsHandler
	HealthcheckHandler *HealthcheckHandler
	logger             logrus.FieldLogger
}

// ErrorHandlingMiddleware handles errors from handlers that return errors
func ErrorHandlingMiddleware(next func(http.ResponseWriter, *http.Request) error) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := next(w, r)
		if err != nil {
			logger := r.Context().Value(utils.LoggerCtxKey)
			if logger != nil {
				logger.(logrus.FieldLogger).Error(err)
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}

func NewHandlers(services *services.Services, logger logrus.FieldLogger) *Handlers {
	validator := validator.New()
	worldsHandler := NewWorldsHandler(services, validator)
	healthcheckHandler := NewHealthcheckHandler()
	return &Handlers{
		logger:             logger,
		WorldsHandler:      worldsHandler,
		HealthcheckHandler: healthcheckHandler,
	}
}

func (h *Handlers) RegisterRoutes(r *mux.Router) {
	logger := h.logger.WithField("method", "RegisterRoutes")
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
	logger := h.logger.WithField("method", "RegisterAuthenticatedRoutes")
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
