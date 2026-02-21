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

func (h *Handler) CreateFromLecture(w http.ResponseWriter, r *http.Request) {
	var req CreateAudioRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("failed to decode request", zap.Error(err))
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.LectureID == uuid.Nil {
		respondWithError(w, http.StatusBadRequest, "lecture_id is required")
		return
	}

	audio, err := h.service.CreateFromLecture(r.Context(), req)
	if err != nil {
		h.logger.Error("failed to create audio", zap.Error(err), zap.String("lecture_id", req.LectureID.String()))
		respondWithError(w, http.StatusInternalServerError, "failed to create audio, please try again")
		return
	}

	respondWithJSON(w, http.StatusCreated, audio)
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
		h.logger.Error("failed to get audio", zap.Error(err), zap.String("audio_id", idStr))
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
		h.logger.Error("failed to get audios", zap.Error(err), zap.String("lecture_id", lectureIDStr))
		respondWithError(w, http.StatusInternalServerError, "failed to retrieve audios")
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
