package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/guilhermeCoutinho/worlds-api/models"
	"github.com/guilhermeCoutinho/worlds-api/services"
)

type WorldsHandler struct {
	services  *services.Services
	validator *validator.Validate
}

func NewWorldsHandler(services *services.Services, validator *validator.Validate) *WorldsHandler {
	return &WorldsHandler{services: services, validator: validator}
}

func (h *WorldsHandler) RegisterHandler(r *mux.Router) {
	r.Handle("/worlds", ErrorHandlingMiddleware(h.HandleGetWorlds)).Methods("GET")
}

func (h *WorldsHandler) RegisterAuthenticatedHandler(r *mux.Router) {
	r.Handle("/worlds", ErrorHandlingMiddleware(h.HandleCreateWorld)).Methods("POST")
	r.Handle("/worlds/{id}", ErrorHandlingMiddleware(h.HandleGetWorldByID)).Methods("GET")
	r.Handle("/worlds/{id}", ErrorHandlingMiddleware(h.HandleUpdateWorld)).Methods("PUT")
}

type GetWorldsQueryParams struct {
	OwnerID string `validate:"omitempty,uuid"`
}

func (h *WorldsHandler) HandleGetWorlds(w http.ResponseWriter, r *http.Request) error {
	params := GetWorldsQueryParams{
		OwnerID: r.URL.Query().Get("ownerId"),
	}
	if err := h.validator.Struct(params); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil
	}

	var worlds []models.World
	var err error

	if params.OwnerID != "" {
		ownerIDUUID := uuid.MustParse(params.OwnerID)
		worlds, err = h.services.WorldsService.GetWorldsByOwnerID(ownerIDUUID)
	} else {
		worlds, err = h.services.WorldsService.GetWorlds()
	}

	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(worlds)
}

type WorldIDParam struct {
	ID string `validate:"required,uuid"`
}

func (h *WorldsHandler) HandleGetWorldByID(w http.ResponseWriter, r *http.Request) error {
	params := WorldIDParam{
		ID: mux.Vars(r)["id"],
	}
	if err := h.validator.Struct(params); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil
	}

	worldID := uuid.MustParse(params.ID)
	world, err := h.services.WorldsService.GetWorldByID(worldID)
	if err != nil {
		http.Error(w, "World not found", http.StatusNotFound)
		return nil
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(world)
}

type CreateWorldRequest struct {
	Name        string `json:"name" validate:"required,max=255,min=3"`
	Description string `json:"description" validate:"required,max=1000,min=3"`
}

func (h *WorldsHandler) HandleCreateWorld(w http.ResponseWriter, r *http.Request) error {
	var req CreateWorldRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil
	}

	if err := h.validator.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil
	}

	userID, err := UserIDFromCtx(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return nil
	}

	world, err := h.services.WorldsService.CreateWorld(userID, req.Name, req.Description)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	return json.NewEncoder(w).Encode(world)
}

type UpdateWorldRequest struct {
	Name        string `json:"name" validate:"required,max=255,min=3"`
	Description string `json:"description" validate:"required,max=1000,min=3"`
}

func (h *WorldsHandler) HandleUpdateWorld(w http.ResponseWriter, r *http.Request) error {
	params := WorldIDParam{
		ID: mux.Vars(r)["id"],
	}
	if err := h.validator.Struct(params); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil
	}

	var req UpdateWorldRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil
	}

	if err := h.validator.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil
	}

	userID, err := UserIDFromCtx(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return nil
	}
	worldID := uuid.MustParse(params.ID)
	world, err := h.services.WorldsService.UpdateWorld(userID, worldID, req.Name, req.Description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return nil
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(world)
}
