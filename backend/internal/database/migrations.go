package database

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RunMigrations(ctx context.Context, pool *pgxpool.Pool, migrationsDir string) error {
	_, err := pool.Exec(ctx, `
                CREATE TABLE IF NOT EXISTS schema_migrations (
                        version BIGINT PRIMARY KEY,
                        name VARCHAR(255) NOT NULL,
                        applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
                )
        `)
	if err != nil {
		return fmt.Errorf("initialize schema migrations table: %w", err)
	}
	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("read migrations directory: %w", err)
	}

	var migrations []string

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()

		if strings.HasSuffix(name, ".up.sql") {
			migrations = append(migrations, name)
		}
	}

	sort.Strings(migrations)

	for _, name := range migrations {
		if err := applyMigration(ctx, pool, migrationsDir, name); err != nil {
			return err
		}
	}

	return nil
}

func applyMigration(ctx context.Context, pool *pgxpool.Pool, migrationsDir, name string) error {
	baseName := strings.TrimSuffix(name, ".up.sql")
	parts := strings.SplitN(baseName, "_", 2)

	if len(parts) != 2 {
		return fmt.Errorf("invalid migration filename: %s", name)
	}

	version, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return fmt.Errorf("parse migration version %s: %w", name, err)
	}

	migrationName := parts[1]

	tx, err := pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin migration transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var applied bool
	err = tx.QueryRow(
		ctx,
		"SELECT EXISTS (SELECT 1 FROM schema_migrations WHERE version = $1)",
		version,
	).Scan(&applied)
	if err != nil {
		return fmt.Errorf("check migration %d: %w", version, err)
	}

	if applied {
		return nil
	}

	sql, err := os.ReadFile(filepath.Join(migrationsDir, name))
	if err != nil {
		return fmt.Errorf("read migration %s: %w", name, err)
	}

	if _, err := tx.Exec(ctx, string(sql)); err != nil {
		return fmt.Errorf("execute migration %s: %w", name, err)
	}

	_, err = tx.Exec(
		ctx,
		"INSERT INTO schema_migrations (version, name) VALUES ($1, $2)",
		version,
		migrationName,
	)
	if err != nil {
		return fmt.Errorf("record migration %s: %w", name, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit migration %s: %w", name, err)
	}

	return nil
}
