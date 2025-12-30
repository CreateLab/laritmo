package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/CreateLab/laritmo/internal/config"
	"github.com/CreateLab/laritmo/internal/database"
	sq "github.com/Masterminds/squirrel"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	ctx := context.Background()

	
	cfg, err := config.Load("configs/config.local.yaml")
	if err != nil {
		slog.ErrorContext(ctx, "Failed to load config", "error", err)
		os.Exit(1)
	}

	
	db, err := database.Connect(cfg.Database.DSN())
	if err != nil {
		slog.ErrorContext(ctx, "Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	var username, email, password string

	
	fmt.Print("Username: ")
	fmt.Scanln(&username)

	fmt.Print("Email: ")
	fmt.Scanln(&email)

	fmt.Print("Password: ")
	fmt.Scanln(&password)

	if username == "" || email == "" || password == "" {
		slog.ErrorContext(ctx, "All fields are required")
		os.Exit(1)
	}

	
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		slog.ErrorContext(ctx, "Password hashing error", "error", err)
		os.Exit(1)
	}

	
	query, args, _ := sq.Insert("users").
		Columns("username", "email", "password_hash", "role").
		Values(username, email, string(passwordHash), "admin").
		ToSql()

	_, err = db.ExecContext(ctx, query, args...)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to create admin", "error", err)
		os.Exit(1)
	}

	slog.InfoContext(ctx, "✅ Admin created successfully!", "username", username, "email", email)
}
