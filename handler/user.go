package handler

import (
	"net/http"

	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/guilhermeCoutinho/worlds-api/models"
	"github.com/guilhermeCoutinho/worlds-api/services"
)

type UserHandler struct {
	services  *services.Services
	validator *validator.Validate
}

func NewUserHandler(services *services.Services, validator *validator.Validate) *UserHandler {
	return &UserHandler{services: services, validator: validator}
}

func (h *UserHandler) RegisterHandler(r *mux.Router) {
	r.Handle("/user/{id}", ErrorHandlingMiddleware(h.HandlerCreateUser)).Methods("POST")
}

type CreateUserParams struct {
	ID string `validate:"required,uuid"`
}

func (h *UserHandler) HandlerCreateUser(w http.ResponseWriter, r *http.Request) error {
	params := CreateUserParams{
		ID: mux.Vars(r)["id"],
	}
	if err := h.validator.Struct(params); err != nil {
		return err
	}
	err := h.services.UserService.CreateUser(&models.User{ID: uuid.MustParse(params.ID)})
	if err != nil {
		return err
	}
	return nil
}
