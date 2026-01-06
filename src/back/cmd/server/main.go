package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/CreateLab/laritmo/docs"
	"github.com/CreateLab/laritmo/internal/auth"
	"github.com/CreateLab/laritmo/internal/config"
	"github.com/CreateLab/laritmo/internal/database"
	"github.com/CreateLab/laritmo/internal/handlers"
	"github.com/CreateLab/laritmo/internal/middleware"
	"github.com/CreateLab/laritmo/internal/repository"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Laritmo API
// @version         1.0
// @description     Educational portal API with course management
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@laritmo.local

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8443
// @BasePath  /
// @schemes   https http

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	ctx := context.Background()

	configPath := os.Getenv("CONFIG_PATH")

	if configPath == "" {
		configPath = "configs/config.local.yaml"
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to load config", "error", err)
		os.Exit(1)
	}

	slog.InfoContext(ctx, "Config loaded successfully")

	db, err := database.Connect(cfg.Database.DSN())
	if err != nil {
		slog.ErrorContext(ctx, "Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()
	courseRepo := repository.NewCourseRepository(db)
	lectureRepo := repository.NewLectureRepository(db)
	labRepo := repository.NewLabRepository(db)
	gradeSheetRepo := repository.NewGradeSheetRepository(db)
	examQuestionRepo := repository.NewExamQuestionRepository(db)
	userRepo := repository.NewUserRepository(db)

	jwtManager := auth.NewJWTManager(cfg.Auth.JWTSecret, cfg.Auth.JWTExpirationHours)
	courseHandler := handlers.NewCourseHandler(courseRepo, logger)
	lectureHandler := handlers.NewLectureHandler(lectureRepo, logger)
	labHandler := handlers.NewLabHandler(labRepo, logger)
	gradeSheetHandler := handlers.NewGradeSheetHandler(gradeSheetRepo, logger)
	examQuestionHandler := handlers.NewExamQuestionHandler(examQuestionRepo, logger)
	authHandler := handlers.NewAuthHandler(userRepo, jwtManager, logger)

	gin.SetMode(cfg.Server.Mode)
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "https://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := r.Group("/api")

	api.GET("/courses", courseHandler.GetAll)
	api.GET("/courses/:id", courseHandler.GetByID)

	api.GET("/lectures", lectureHandler.GetAll)
	api.GET("/lectures/:id", lectureHandler.GetByID)

	api.GET("/labs", labHandler.GetAll)
	api.GET("/labs/:id", labHandler.GetByID)

	api.GET("/grade-sheets", gradeSheetHandler.GetAll)
	api.GET("/grade-sheets/:id", gradeSheetHandler.GetByID)

	api.GET("/exam-questions", examQuestionHandler.GetAll)
	api.GET("/exam-questions/:id", examQuestionHandler.GetByID)

	loginGroup := api.Group("/auth")
	loginGroup.Use(middleware.RateLimitMiddleware(cfg.Auth.GetRateLimitRequests(), cfg.Auth.GetRateLimitBurst()))
	loginGroup.POST("/login", authHandler.Login)

	admin := r.Group("/api/admin")
	admin.Use(middleware.AuthMiddleware(jwtManager))
	admin.Use(middleware.AdminOnly())
	{
		admin.POST("/courses", courseHandler.Create)
		admin.PUT("/courses/:id", courseHandler.Update)
		admin.DELETE("/courses/:id", courseHandler.Delete)

		admin.POST("/lectures", lectureHandler.Create)
		admin.PUT("/lectures/:id", lectureHandler.Update)
		admin.DELETE("/lectures/:id", lectureHandler.Delete)

		admin.POST("/labs", labHandler.Create)
		admin.PUT("/labs/:id", labHandler.Update)
		admin.DELETE("/labs/:id", labHandler.Delete)

		admin.POST("/grade-sheets", gradeSheetHandler.Create)
		admin.PUT("/grade-sheets/:id", gradeSheetHandler.Update)
		admin.DELETE("/grade-sheets/:id", gradeSheetHandler.Delete)

		admin.POST("/exam-questions", examQuestionHandler.Create)
		admin.POST("/exam-questions/bulk", examQuestionHandler.BulkCreateJSON)
		admin.POST("/exam-questions/upload", examQuestionHandler.BulkUploadFile)
		admin.PUT("/exam-questions/:id", examQuestionHandler.Update)
		admin.DELETE("/exam-questions/:id", examQuestionHandler.Delete)
	}

	r.Static("/assets", "./web/assets")
	r.StaticFile("/favicon.ico", "./web/favicon.ico")

	r.NoRoute(func(c *gin.Context) {
		c.File("./web/index.html")
	})

	if cfg.Server.Mode == "debug" {
		protocol := "http"
		if cfg.Server.UseTLS {
			protocol = "https"
		}
		docURL := fmt.Sprintf("%s://%s:%d/swagger/doc.json", protocol, cfg.Server.Host, cfg.Server.Port)

		r.GET("/swagger/*any", ginSwagger.WrapHandler(
			swaggerFiles.Handler,
			ginSwagger.URL(docURL),
		))

		swaggerURL := fmt.Sprintf("%s://%s:%d/swagger/index.html", protocol, cfg.Server.Host, cfg.Server.Port)
		slog.InfoContext(ctx, "Swagger UI available", "url", swaggerURL)
	}

	addr := fmt.Sprintf(":%d", cfg.Server.Port)

	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	go func() {
		protocol := "HTTP"
		url := fmt.Sprintf("http://%s:%d", cfg.Server.Host, cfg.Server.Port)

		var err error
		if cfg.Server.UseTLS {
			protocol = "HTTPS"
			url = fmt.Sprintf("https://%s:%d", cfg.Server.Host, cfg.Server.Port)
			slog.InfoContext(ctx, fmt.Sprintf("%s server started", protocol), "url", url)
			err = srv.ListenAndServeTLS(cfg.Server.TLSCertFile, cfg.Server.TLSKeyFile)
		} else {
			slog.InfoContext(ctx, fmt.Sprintf("%s server started", protocol), "url", url)
			err = srv.ListenAndServe()
		}

		if err != nil && err != http.ErrServerClosed {
			slog.ErrorContext(ctx, "Server startup error", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.InfoContext(ctx, "Shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.ErrorContext(ctx, "Server shutdown error", "error", err)
		os.Exit(1)
	}

	slog.InfoContext(ctx, "Server stopped successfully")
}
