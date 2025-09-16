package handler

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
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
}

func (h *WorldsHandler) HandleGetWorlds(w http.ResponseWriter, r *http.Request) error {
	worlds, err := h.services.WorldsService.GetWorlds()
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Found %d worlds", len(worlds))))
	return nil
}
