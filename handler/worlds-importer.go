package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/guilhermeCoutinho/worlds-api/services"
)

type WorldsImporterHandler struct {
	services  *services.Services
	validator *validator.Validate
}

func NewWorldsImporterHandler(services *services.Services, validator *validator.Validate) *WorldsImporterHandler {
	return &WorldsImporterHandler{services: services, validator: validator}
}

func (h *WorldsImporterHandler) RegisterHandler(r *mux.Router) {
	r.Handle("/worlds/import", ErrorHandlingMiddleware(h.HandleImportWorlds)).Methods("POST")
	r.Handle("/jobs/status/{id}", ErrorHandlingMiddleware(h.HandleGetJobStatus)).Methods("GET")
}

type ImportWorldsRequest struct {
	Worlds            []uuid.UUID `json:"worlds"`
	TargetEnvironment string      `json:"target_environment"`
}

func (h *WorldsImporterHandler) HandleImportWorlds(w http.ResponseWriter, r *http.Request) error {
	var req ImportWorldsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}

	if err := h.validator.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}

	userID, err := UserIDFromCtx(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return err
	}

	response, err := h.services.WorldsImporterService.CreateImportWorldsJob(r.Context(), userID, req.Worlds, req.TargetEnvironment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(response)
}

type GetJobStatusRequest struct {
	ID string `json:"id"`
}

func (h *WorldsImporterHandler) HandleGetJobStatus(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["id"]
	job, err := h.services.WorldsImporterService.GetAndUpdateWorldsTransferJobStatus(r.Context(), uuid.MustParse(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(job)
}
