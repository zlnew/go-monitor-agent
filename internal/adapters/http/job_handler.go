package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"horizonx-server/internal/domain"
)

type JobHandler struct {
	svc domain.JobService
}

func NewJobHandler(svc domain.JobService) *JobHandler {
	return &JobHandler{svc: svc}
}

func (h *JobHandler) Index(w http.ResponseWriter, r *http.Request) {
	jobs, err := h.svc.Get(r.Context())
	if err != nil {
		JSONError(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	JSONSuccess(w, http.StatusOK, APIResponse{
		Message: "OK",
		Data:    jobs,
	})
}

func (h *JobHandler) Show(w http.ResponseWriter, r *http.Request) {
	jobID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		JSONError(w, http.StatusBadRequest, "invalid application id")
		return
	}

	job, err := h.svc.GetByID(r.Context(), jobID)
	if err != nil {
		if errors.Is(err, domain.ErrJobNotFound) {
			JSONError(w, http.StatusNotFound, "job not found")
			return
		}
		JSONError(w, http.StatusInternalServerError, "failed to get job")
		return
	}

	JSONSuccess(w, http.StatusOK, APIResponse{
		Message: "OK",
		Data:    job,
	})
}

func (h *JobHandler) Start(w http.ResponseWriter, r *http.Request) {
	paramID := r.PathValue("id")

	jobID, err := strconv.ParseInt(paramID, 10, 64)
	if err != nil {
		JSONError(w, http.StatusBadRequest, "invalid job id")
		return
	}

	job, err := h.svc.Start(r.Context(), jobID)
	if err != nil {
		if errors.Is(err, domain.ErrJobNotFound) {
			JSONError(w, http.StatusNotFound, "job not found")
			return
		}

		JSONError(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	JSONSuccess(w, http.StatusOK, APIResponse{
		Message: "OK",
		Data:    job,
	})
}

func (h *JobHandler) Finish(w http.ResponseWriter, r *http.Request) {
	paramID := r.PathValue("id")

	var req domain.JobFinishRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		JSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if validationErrors := ValidateStruct(req); len(validationErrors) > 0 {
		JSONValidationError(w, validationErrors)
		return
	}

	jobID, err := strconv.ParseInt(paramID, 10, 64)
	if err != nil {
		JSONError(w, http.StatusBadRequest, "invalid job id")
		return
	}

	job, err := h.svc.Finish(r.Context(), jobID, req.Status, &req.OutputLog)
	if err != nil {
		if errors.Is(err, domain.ErrJobNotFound) {
			JSONError(w, http.StatusNotFound, "job not found")
			return
		}

		log.Println("asd", err.Error())

		JSONError(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	JSONSuccess(w, http.StatusOK, APIResponse{
		Message: "OK",
		Data:    job,
	})
}
