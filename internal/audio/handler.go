package audio

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Handler struct {
	service Service
	logger  *zap.Logger
}

func NewHandler(service Service, logger *zap.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid audio id")
		return
	}

	audio, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		h.logger.Error("failed to get audio", zap.Error(err))
		respondWithError(w, http.StatusNotFound, "audio not found")
		return
	}

	respondWithJSON(w, http.StatusOK, audio)
}

func (h *Handler) GetByLectureID(w http.ResponseWriter, r *http.Request) {
	lectureIDStr := chi.URLParam(r, "lecture_id")
	lectureID, err := uuid.Parse(lectureIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid lecture id")
		return
	}

	audios, err := h.service.GetByLectureID(r.Context(), lectureID)
	if err != nil {
		h.logger.Error("failed to get audios", zap.Error(err))
		respondWithError(w, http.StatusInternalServerError, "failed to get audios")
		return
	}

	respondWithJSON(w, http.StatusOK, audios)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}
