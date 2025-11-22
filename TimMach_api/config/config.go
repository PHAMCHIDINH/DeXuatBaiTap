package config

import (
	"strings"
	"time"

	"github.com/caarlos0/env/v10"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Config struct {
	DBURL     string `env:"DB_URL" envDefault:"postgres://postgres:postgres@localhost:5432/heartdb?sslmode=disable"`
	JWTSecret string `env:"JWT_SECRET" envDefault:"dev-secret"`
	MLBaseURL string `env:"ML_BASE_URL" envDefault:"http://localhost:8000"`
	Port      string `env:"PORT" envDefault:"8080"`
}

// Load đọc biến môi trường vào Config.
func Load() (Config, error) {
	var cfg Config
	err := env.Parse(&cfg)
	return cfg, err
}

// CORSMiddleware tạo middleware CORS theo cấu hình.
func (c Config) CORSMiddleware() gin.HandlerFunc {
	cfg := cors.DefaultConfig()
	cfg.AllowOrigins = []string{"http://localhost:5173"}
	cfg.AllowCredentials = true
	cfg.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	cfg.AllowHeaders = []string{"Authorization", "Content-Type"}
	cfg.MaxAge = 12 * time.Hour
	return cors.New(cfg)
}

func splitAndTrim(input string) []string {
	parts := strings.Split(input, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
