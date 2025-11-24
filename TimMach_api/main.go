package main

import (
	"context"
	"fmt"
	"time"

	"chidinh/config"
	db "chidinh/db/sqlc"
	"chidinh/middleware"
	"chidinh/modules/auth"
	"chidinh/modules/exercises"
	"chidinh/modules/patients"
	"chidinh/modules/predictions"
	"chidinh/modules/reports"
	"chidinh/modules/stats"
	"chidinh/modules/users"
	"chidinh/utils"
	"chidinh/utils/httpclient"
	"chidinh/utils/mailer"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env for local dev; ignore error if file is missing.
	_ = godotenv.Load(".env")

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

	mlHTTPClient := httpclient.NewRestyClient(cfg.MLBaseURL, 10*time.Second)
	mailerSvc := mailer.New(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPass, cfg.SMTPFrom)

	// Handlers
	userController := users.NewController(queries, tokenMaker)
	patientController := patients.NewController(queries)
	reportController := reports.NewController(queries, mailerSvc, cfg.ReportDefaultEmail)
	predictionController := predictions.NewController(queries, mlHTTPClient)
	exerciseController := exercises.NewController(queries)
	statsController := stats.NewController(queries)

	// Router
	router := gin.Default()
	router.Use(cfg.CORSMiddleware())
	api := router.Group("/api")
	authMiddleware := middleware.AuthMiddleware(cfg.JWTSecret)

	users.RegisterUserRoutes(api, userController, authMiddleware)
	patients.RegisterPatientRoutes(api, patientController, authMiddleware)
	predictions.RegisterPredictionRoutes(api, predictionController, authMiddleware)
	exercises.RegisterExerciseRoutes(api, exerciseController, authMiddleware)
	reports.RegisterReportRoutes(api, reportController, authMiddleware)
	stats.RegisterStatsRoutes(api, statsController, authMiddleware)

	if err := router.Run(":" + cfg.Port); err != nil {
		logger.Fatalw("server exited", "error", err)
	}
}
