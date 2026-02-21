package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	App     AppConfig
	DB      DBConfig
	JWT     JWTConfig
	MinIO   MinIOConfig
	API     APIConfig
}

type AppConfig struct {
	Port string
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type JWTConfig struct {
	Secret     string
	Expiration time.Duration
}

type MinIOConfig struct {
	Endpoint      string
	AccessKey     string
	SecretKey     string
	BucketAudios  string
	BucketVideos  string
	BucketAvatars string
	UseSSL        bool
}

type APIConfig struct {
	GeminiKey    string
	GeminiModel  string
	VoiceRSSKey  string
	SyncKey      string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	jwtExp, err := time.ParseDuration(getEnv("JWT_EXPIRATION", "24h"))
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_EXPIRATION: %w", err)
	}

	cfg := &Config{
		App: AppConfig{
			Port: getEnv("APP_PORT", "8080"),
		},
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			Name:     getEnv("DB_NAME", "mlp"),
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", "supersecretkey"),
			Expiration: jwtExp,
		},
		MinIO: MinIOConfig{
			Endpoint:      getEnv("MINIO_ENDPOINT", "localhost:9000"),
			AccessKey:     getEnv("MINIO_ACCESS_KEY", "minioadmin"),
			SecretKey:     getEnv("MINIO_SECRET_KEY", "minioadmin"),
			BucketAudios:  getEnv("MINIO_BUCKET_AUDIOS", "audios"),
			BucketVideos:  getEnv("MINIO_BUCKET_VIDEOS", "videos"),
			BucketAvatars: getEnv("MINIO_BUCKET_AVATARS", "avatars"),
			UseSSL:        getEnv("MINIO_USE_SSL", "false") == "true",
		},
		API: APIConfig{
			GeminiKey:   getEnv("GEMINI_API_KEY", ""),
			GeminiModel: getEnv("GEMINI_MODEL", "gemini-2.0-flash-lite"),
			VoiceRSSKey: getEnv("VOICERSS_API_KEY", ""),
			SyncKey:     getEnv("SYNC_API_KEY", ""),
		},
	}

	return cfg, nil
}

func (c *DBConfig) ConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.User, c.Password, c.Host, c.Port, c.Name)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
