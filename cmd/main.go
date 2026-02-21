package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"mlp/internal/audio"
	"mlp/internal/config"
	"mlp/internal/lecture"
	"mlp/internal/middleware"
	"mlp/internal/storage"
	"mlp/internal/user"
	"mlp/internal/video"
	"mlp/pkg/gemini"
	"mlp/pkg/lipsync"
	"mlp/pkg/voicerss"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	logger.Info("starting mlp server")

	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("failed to load config", zap.Error(err))
	}

	ctx := context.Background()

	dbPool, err := pgxpool.New(ctx, cfg.DB.ConnectionString())
	if err != nil {
		logger.Fatal("failed to connect to database", zap.Error(err))
	}
	defer dbPool.Close()

	if err := dbPool.Ping(ctx); err != nil {
		logger.Fatal("failed to ping database", zap.Error(err))
	}
	logger.Info("connected to database")

	minioClient, err := storage.NewMinIOClient(
		cfg.MinIO.Endpoint,
		cfg.MinIO.AccessKey,
		cfg.MinIO.SecretKey,
		cfg.MinIO.BucketAudios,
		cfg.MinIO.BucketVideos,
		cfg.MinIO.BucketAvatars,
		cfg.MinIO.UseSSL,
		logger,
	)
	if err != nil {
		logger.Fatal("failed to create minio client", zap.Error(err))
	}

	if err := minioClient.EnsureBuckets(ctx); err != nil {
		logger.Fatal("failed to ensure minio buckets", zap.Error(err))
	}
	logger.Info("connected to minio and ensured buckets")

	geminiClient := gemini.NewClient(cfg.API.CohereKey, cfg.API.CohereModel)
	voiceRSSClient := voicerss.NewClient(cfg.API.VoiceRSSKey)
	lipsyncClient := lipsync.NewClient("http://lipsync:5000")

	userRepo := user.NewRepository(dbPool)
	userService := user.NewService(userRepo, cfg.JWT.Secret, cfg.JWT.Expiration)
	userHandler := user.NewHandler(userService, logger)

	lectureRepo := lecture.NewRepository(dbPool)
	lectureService := lecture.NewService(lectureRepo, geminiClient)
	lectureHandler := lecture.NewHandler(lectureService, logger)

	audioRepo := audio.NewRepository(dbPool)
	audioService := audio.NewService(audioRepo, lectureRepo, voiceRSSClient, minioClient, cfg.MinIO.BucketAudios)
	audioHandler := audio.NewHandler(audioService, logger)

	videoRepo := video.NewRepository(dbPool)
	videoService := video.NewService(videoRepo, audioRepo, minioClient, lipsyncClient, cfg.MinIO.BucketVideos, cfg.MinIO.BucketAvatars, cfg.MinIO.Endpoint)
	videoHandler := video.NewHandler(videoService, logger)

	authMiddleware := middleware.NewAuthMiddleware(cfg.JWT.Secret, logger)

	r := chi.NewRouter()

	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	r.Route("/api/v1/auth", func(r chi.Router) {
		r.Post("/register", userHandler.Register)
		r.Post("/login", userHandler.Login)
	})

	r.Route("/api/v1/lectures", func(r chi.Router) {
		r.Use(authMiddleware.Authenticate)
		r.Post("/", lectureHandler.GenerateLecture)
		r.Get("/", lectureHandler.GetByUserID)
		r.Get("/{id}", lectureHandler.GetByID)
	})

	r.Route("/api/v1/audios", func(r chi.Router) {
		r.Use(authMiddleware.Authenticate)
		r.Post("/", audioHandler.CreateFromLecture)
		r.Get("/{id}", audioHandler.GetByID)
		r.Get("/lecture/{lecture_id}", audioHandler.GetByLectureID)
	})

	r.Route("/api/v1/videos", func(r chi.Router) {
		r.Use(authMiddleware.Authenticate)
		r.Post("/", videoHandler.CreateFromAudio)
		r.Get("/{id}", videoHandler.GetByID)
		r.Get("/audio/{audio_id}", videoHandler.GetByAudioID)
	})

	server := &http.Server{
		Addr:         ":" + cfg.App.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		logger.Info("server listening", zap.String("port", cfg.App.Port))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("failed to start server", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("server forced to shutdown", zap.Error(err))
	}

	logger.Info("server exited")
}
