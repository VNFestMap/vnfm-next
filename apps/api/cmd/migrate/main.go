package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"vnfm-api/internal/infrastructure/database"
	"vnfm-api/pkg/config"
	"vnfm-api/pkg/logger"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	direction := flag.String("dir", "up", "Migration direction: up or down")
	step := flag.Int("step", 0, "Number of migrations to run (0 = all)")
	migrationsDir := flag.String("path", "migrations", "Path to migration files")
	flag.Parse()

	cfg, err := config.Load()
	if err != nil {
		slog.Error("load config", "error", err)
		os.Exit(1)
	}
	if err := config.RequireDatabase(cfg); err != nil {
		slog.Error("config", "error", err)
		os.Exit(1)
	}
	logger.Init(cfg.Server.Mode)

	db, err := database.Open(cfg.Database, cfg.Server.Mode)
	if err != nil {
		slog.Error("open database", "error", err)
		os.Exit(1)
	}
	sqlDB, err := db.DB()
	if err != nil {
		slog.Error("sql db", "error", err)
		os.Exit(1)
	}

	_, err = sqlDB.Exec(`
		CREATE TABLE IF NOT EXISTS _migrations (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL UNIQUE,
			applied_at TIMESTAMP NOT NULL DEFAULT NOW()
		)
	`)
	if err != nil {
		slog.Error("create _migrations", "error", err)
		os.Exit(1)
	}

	rows, err := sqlDB.Query("SELECT name FROM _migrations ORDER BY id")
	if err != nil {
		slog.Error("list migrations", "error", err)
		os.Exit(1)
	}
	defer rows.Close()

	applied := map[string]bool{}
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			slog.Error("scan migration", "error", err)
			os.Exit(1)
		}
		applied[name] = true
	}

	suffix := "." + *direction + ".sql"
	files, err := filepath.Glob(filepath.Join(*migrationsDir, "*"+suffix))
	if err != nil {
		slog.Error("glob migrations", "error", err)
		os.Exit(1)
	}
	sort.Strings(files)
	if *direction == "down" {
		for i, j := 0, len(files)-1; i < j; i, j = i+1, j-1 {
			files[i], files[j] = files[j], files[i]
		}
	}

	ran := 0
	for _, file := range files {
		base := filepath.Base(file)
		name := strings.TrimSuffix(base, suffix)

		if *direction == "up" {
			if applied[name] {
				continue
			}
		} else if !applied[name] {
			continue
		}

		if *step > 0 && ran >= *step {
			break
		}

		slog.Info("running migration", "file", base, "direction", *direction)
		content, err := os.ReadFile(file)
		if err != nil {
			slog.Error("read migration", "file", base, "error", err)
			os.Exit(1)
		}
		if _, err = sqlDB.Exec(string(content)); err != nil {
			slog.Error("exec migration", "file", base, "error", err)
			os.Exit(1)
		}
		if *direction == "up" {
			_, err = sqlDB.Exec("INSERT INTO _migrations (name) VALUES ($1)", name)
		} else {
			_, err = sqlDB.Exec("DELETE FROM _migrations WHERE name = $1", name)
		}
		if err != nil {
			slog.Error("record migration", "file", base, "error", err)
			os.Exit(1)
		}
		ran++
	}

	if ran == 0 {
		fmt.Println("没有待执行的迁移")
	} else {
		fmt.Printf("成功执行 %d 个迁移\n", ran)
	}
}
