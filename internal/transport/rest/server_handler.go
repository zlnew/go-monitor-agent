package rest

import (
	"encoding/json"
	"net/http"

	"horizonx-server/internal/domain"
)

type ServerHandler struct {
	svc domain.ServerService
}

func NewServerHandler(svc domain.ServerService) *ServerHandler {
	return &ServerHandler{svc: svc}
}

func (h *ServerHandler) Index(w http.ResponseWriter, r *http.Request) {
	servers, err := h.svc.List(r.Context())
	if err != nil {
		JSONError(w, http.StatusInternalServerError, "Failed to retrieve server data")
		return
	}

	JSONSuccess(w, http.StatusOK, APIResponse{
		Message: "OK",
		Data:    servers,
	})
}

func (h *ServerHandler) Store(w http.ResponseWriter, r *http.Request) {
	var req domain.ServerCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		JSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if validationErrors := ValidateStruct(req); len(validationErrors) > 0 {
		JSONValidationError(w, validationErrors)
		return
	}

	srv, token, err := h.svc.Register(r.Context(), req.Name, req.IPAddress)
	if err != nil {
		JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	JSONSuccess(w, http.StatusCreated, APIResponse{
		Message: "Server registered. COPY THIS TOKEN NOW!",
		Data: map[string]any{
			"server": srv,
			"token":  token,
		},
	})
}
