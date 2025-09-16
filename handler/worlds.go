package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/guilhermeCoutinho/worlds-api/models"
	"github.com/guilhermeCoutinho/worlds-api/services"
)

type WorldsHandler struct {
	services *services.Services
}

func NewWorldsHandler(services *services.Services) *WorldsHandler {
	return &WorldsHandler{services: services}
}

func (h *WorldsHandler) RegisterAuthenticatedHandler(r *mux.Router) {
	r.HandleFunc("/worlds", LoggedHandler(h.HandleGetWorlds)).Methods("GET")
	r.HandleFunc("/worlds", LoggedHandler(h.HandleCreateWorld)).Methods("POST")
	r.HandleFunc("/worlds/{id}", LoggedHandler(h.HandleGetWorldByID)).Methods("GET")
	r.HandleFunc("/worlds/{id}", LoggedHandler(h.HandleUpdateWorld)).Methods("PUT")
}

type CreateWorldRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type UpdateWorldRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
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
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return nil
	}

	if req.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
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
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return nil
	}

	if req.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return nil
	}

	world, err := h.services.WorldsService.UpdateWorld(id, req.Name, req.Description)
	if err != nil {
		http.Error(w, "World not found or update failed", http.StatusNotFound)
		return nil
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(world)
}
