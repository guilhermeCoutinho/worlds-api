package handler

import (
	"net/http"

	"github.com/gorilla/mux"
)

type HealthcheckHandler struct {
}

func NewHealthcheckHandler() *HealthcheckHandler {
	return &HealthcheckHandler{}
}

func (h *HealthcheckHandler) RegisterHandler(r *mux.Router) {
	r.HandleFunc("/healthcheck", LoggedHandler(h.HandleHealthcheck)).Methods("GET", "OPTIONS")
}

func (h *HealthcheckHandler) HandleHealthcheck(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
	return nil
}
