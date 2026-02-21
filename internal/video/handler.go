package video

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
		respondWithError(w, http.StatusBadRequest, "invalid video id")
		return
	}

	video, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		h.logger.Error("failed to get video", zap.Error(err))
		respondWithError(w, http.StatusNotFound, "video not found")
		return
	}

	respondWithJSON(w, http.StatusOK, video)
}

func (h *Handler) GetByAudioID(w http.ResponseWriter, r *http.Request) {
	audioIDStr := chi.URLParam(r, "audio_id")
	audioID, err := uuid.Parse(audioIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid audio id")
		return
	}

	videos, err := h.service.GetByAudioID(r.Context(), audioID)
	if err != nil {
		h.logger.Error("failed to get videos", zap.Error(err))
		respondWithError(w, http.StatusInternalServerError, "failed to get videos")
		return
	}

	respondWithJSON(w, http.StatusOK, videos)
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
