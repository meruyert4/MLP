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

func (h *Handler) CreateFromAudio(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		h.logger.Error("failed to parse multipart form", zap.Error(err))
		respondWithError(w, http.StatusBadRequest, "invalid form data")
		return
	}

	audioIDStr := r.FormValue("audio_id")
	if audioIDStr == "" {
		respondWithError(w, http.StatusBadRequest, "audio_id is required")
		return
	}

	audioID, err := uuid.Parse(audioIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid audio_id")
		return
	}

	file, header, err := r.FormFile("avatar")
	if err != nil {
		h.logger.Error("failed to get avatar file", zap.Error(err))
		respondWithError(w, http.StatusBadRequest, "avatar file is required")
		return
	}
	defer file.Close()

	req := CreateVideoRequest{
		AudioID: audioID,
	}

	video, err := h.service.CreateFromAudio(r.Context(), req, file, header.Filename)
	if err != nil {
		h.logger.Error("failed to create video", zap.Error(err), zap.String("audio_id", audioIDStr))
		respondWithError(w, http.StatusInternalServerError, "failed to create video, please try again")
		return
	}

	respondWithJSON(w, http.StatusAccepted, video)
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
		h.logger.Error("failed to get video", zap.Error(err), zap.String("video_id", idStr))
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
		h.logger.Error("failed to get videos", zap.Error(err), zap.String("audio_id", audioIDStr))
		respondWithError(w, http.StatusInternalServerError, "failed to retrieve videos")
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
