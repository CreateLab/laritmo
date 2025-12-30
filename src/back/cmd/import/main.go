package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/CreateLab/laritmo/internal/config"
	"github.com/CreateLab/laritmo/internal/database"
	sq "github.com/Masterminds/squirrel"
)

const repoPath = "../../tmp/AspITMO" 

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

	
	courseID, err := createCourse(ctx, db)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to create course", "error", err)
		os.Exit(1)
	}
	slog.InfoContext(ctx, "✅ Course created", "course_id", courseID)

	
	if err := importLectures(ctx, db, courseID); err != nil {
		slog.ErrorContext(ctx, "Failed to import lectures", "error", err)
		os.Exit(1)
	}

	
	if err := importLabs(ctx, db, courseID); err != nil {
		slog.ErrorContext(ctx, "Failed to import labs", "error", err)
		os.Exit(1)
	}

	slog.InfoContext(ctx, "🎉 Import completed successfully!")
}

func createCourse(ctx context.Context, db *sql.DB) (int64, error) {
	query, args, _ := sq.Insert("courses").
		Columns("name", "semester", "description").
		Values("ASP.NET Core", "2024-2025", "Основы разработки Web API на платформе ASP.NET Core").
		ToSql()

	result, err := db.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to create course: %w", err)
	}

	return result.LastInsertId()
}

func importLectures(ctx context.Context, db *sql.DB, courseID int64) error {
	lecturesPath := filepath.Join(repoPath, "lections")

	
	files, err := filepath.Glob(filepath.Join(lecturesPath, "L*.md"))
	if err != nil {
		return err
	}

	slog.InfoContext(ctx, "📚 Found lectures", "count", len(files))

	for _, file := range files {
		
		weekNum := extractNumber(filepath.Base(file), "L")

		
		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", file, err)
		}

		
		title := extractTitle(string(content))

		// GitHub URL
		githubURL := fmt.Sprintf("https://github.com/CreateLab/AspITMO/blob/main/lections/%s", filepath.Base(file))

		
		query, args, _ := sq.Insert("lectures").
			Columns("course_id", "week", "title", "content", "github_url").
			Values(courseID, weekNum, title, string(content), githubURL).
			ToSql()

		_, err = db.ExecContext(ctx, query, args...)
		if err != nil {
			return fmt.Errorf("failed to insert lecture %s: %w", file, err)
		}

		slog.InfoContext(ctx, "  ✓ Lecture imported", "week", weekNum, "title", truncate(title, 50))
	}

	return nil
}

func importLabs(ctx context.Context, db *sql.DB, courseID int64) error {
	
	files, err := filepath.Glob(filepath.Join(repoPath, "Lab*.md"))
	if err != nil {
		return err
	}

	slog.InfoContext(ctx, "🔬 Found labs", "count", len(files))

	for _, file := range files {
		
		labNum := extractNumber(filepath.Base(file), "Lab")

		
		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", file, err)
		}

		
		title := extractTitle(string(content))

		// GitHub URL
		githubURL := fmt.Sprintf("https://github.com/CreateLab/AspITMO/blob/main/%s", filepath.Base(file))

		
		query, args, _ := sq.Insert("labs").
			Columns("course_id", "number", "title", "description", "max_score", "github_url").
			Values(courseID, labNum, title, string(content), 100, githubURL).
			ToSql()

		_, err = db.ExecContext(ctx, query, args...)
		if err != nil {
			return fmt.Errorf("failed to insert lab %s: %w", file, err)
		}

		slog.InfoContext(ctx, "  ✓ Lab imported", "number", labNum, "title", truncate(title, 50))
	}

	return nil
}


func extractTitle(content string) string {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "# ") {
			return strings.TrimPrefix(line, "# ")
		}
	}
	return "Untitled"
}


func extractNumber(filename, prefix string) int {
	re := regexp.MustCompile(prefix + `(\d+)\.md`)
	matches := re.FindStringSubmatch(filename)
	if len(matches) > 1 {
		var num int
		fmt.Sscanf(matches[1], "%d", &num)
		return num
	}
	return 0
}


func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}
