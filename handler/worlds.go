package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator"
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

type CreateWorldRequest struct {
	Name        string `json:"name" validate:"required,max=255,min=3"`
	Description string `json:"description" validate:"required,max=1000,min=3"`
}

type UpdateWorldRequest struct {
	Name        string `json:"name" validate:"required,max=255,min=3"`
	Description string `json:"description" validate:"required,max=1000,min=3"`
}

func (h *WorldsHandler) HandleGetWorlds(w http.ResponseWriter, r *http.Request) error {
	ownerID := r.URL.Query().Get("ownerId")

	var worlds []models.World
	var err error

	if ownerID != "" {
		worlds, err = h.services.WorldsService.GetWorldsByOwnerID(ownerID)
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

func (h *WorldsHandler) HandleGetWorldByID(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id := vars["id"]

	world, err := h.services.WorldsService.GetWorldByID(id)
	if err != nil {
		http.Error(w, "World not found", http.StatusNotFound)
		return nil
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(world)
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

	userID := r.Context().Value("userID").(string)

	world, err := h.services.WorldsService.CreateWorld(req.Name, req.Description, userID)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	return json.NewEncoder(w).Encode(world)
}

func (h *WorldsHandler) HandleUpdateWorld(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id := vars["id"]

	var req UpdateWorldRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil
	}

	if err := h.validator.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil
	}

	world, err := h.services.WorldsService.UpdateWorld(id, req.Name, req.Description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return nil
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(world)
}
