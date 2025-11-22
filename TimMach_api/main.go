package main

import (
	"context"
	"fmt"
	"time"

	"chidinh/config"
	db "chidinh/db/sqlc"
	"chidinh/middleware"
	"chidinh/modules/auth"
	"chidinh/modules/patients"
	"chidinh/modules/predictions"
	"chidinh/modules/exercises"
	"chidinh/modules/stats"
	"chidinh/modules/users"
	"chidinh/utils"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	//"github.com/joho/godotenv"
)

func main() {
	// Load .env for local dev; ignore error if file is missing.
	//_ = godotenv.Load(".env")

	logger := utils.InitLogger()
	defer func() { _ = logger.Sync() }()

	cfg, err := config.Load()
	fmt.Println("JWT Secret:", cfg.JWTSecret)
	fmt.Println("DB URL:", cfg.DBURL)
	if err != nil {
		logger.Fatalw("cannot parse env", "error", err)
	}

	// Init DB
	pool, err := pgxpool.New(context.Background(), cfg.DBURL)
	if err != nil {
		logger.Fatalw("cannot connect db", "error", err, "db_url", cfg.DBURL)
	}
	defer pool.Close()
	queries := db.New(pool)

	// Services
	tokenMaker := auth.JWTMaker{Secret: cfg.JWTSecret, TTL: 24 * time.Hour}

	mlClient := predictions.NewRestyMLClient(cfg.MLBaseURL, 10*time.Second)

	// Handlers
	userHandler := users.NewHandler(queries, tokenMaker)
	patientHandler := patients.NewHandler(queries)
	predictionHandler := predictions.NewHandler(queries, mlClient)
	exerciseHandler := exercises.NewHandler(queries)
	statsHandler := stats.NewHandler(queries)

	// Router
	router := gin.Default()
	router.Use(cfg.CORSMiddleware())
	api := router.Group("/api")
	authMiddleware := middleware.AuthMiddleware(cfg.JWTSecret)

	users.RegisterUserRoutes(api, userHandler, authMiddleware)
	patients.RegisterPatientRoutes(api, patientHandler, authMiddleware)
	predictions.RegisterPredictionRoutes(api, predictionHandler, authMiddleware)
	exercises.RegisterExerciseRoutes(api, exerciseHandler, authMiddleware)
	stats.RegisterStatsRoutes(api, statsHandler, authMiddleware)

	if err := router.Run(":" + cfg.Port); err != nil {
		logger.Fatalw("server exited", "error", err)
	}
}
