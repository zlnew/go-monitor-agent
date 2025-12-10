package rest

import (
	"encoding/json"
	"net/http"

	"horizonx-server/internal/domain"
	"horizonx-server/internal/transport/rest/middleware"
)

type MetricsHandler struct {
	svc domain.MetricsService
}

func NewMetricsHandler(svc domain.MetricsService) *MetricsHandler {
	return &MetricsHandler{svc: svc}
}

func (h *MetricsHandler) Report(w http.ResponseWriter, r *http.Request) {
	serverID, ok := r.Context().Value(middleware.ServerIDKey).(int64)
	if !ok {
		JSONError(w, http.StatusInternalServerError, "Server Context Missing")
		return
	}

	var m domain.Metrics
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		JSONError(w, http.StatusBadRequest, "Bad Request")
		return
	}

	m.ServerID = serverID

	if err := h.svc.Ingest(r.Context(), m); err != nil {
		JSONError(w, http.StatusBadRequest, "Internal Server Error")
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
