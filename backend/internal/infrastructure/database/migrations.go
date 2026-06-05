package database

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RunMigrations(ctx context.Context, db *pgxpool.Pool, migrationsPath string) error {
	entries, err := os.ReadDir(migrationsPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("migrations path not found: %s", migrationsPath)
		}
		return err
	}

	files := make([]string, 0)
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}
		files = append(files, entry.Name())
	}
	sort.Strings(files)
	if len(files) == 0 {
		return nil
	}

	if err := ensureMigrationTable(ctx, db); err != nil {
		return err
	}
	if err := baselineExistingInit(ctx, db, files); err != nil {
		return err
	}

	for _, file := range files {
		applied, err := migrationApplied(ctx, db, file)
		if err != nil {
			return err
		}
		if applied {
			continue
		}
		sqlBytes, err := os.ReadFile(filepath.Join(migrationsPath, file))
		if err != nil {
			return err
		}
		if err := runMigration(ctx, db, file, string(sqlBytes)); err != nil {
			return err
		}
		log.Printf("migration applied: %s", file)
	}
	return nil
}

func ensureMigrationTable(ctx context.Context, db *pgxpool.Pool) error {
	_, err := db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version TEXT PRIMARY KEY,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT now()
		)`)
	return err
}

func baselineExistingInit(ctx context.Context, db *pgxpool.Pool, files []string) error {
	hasApplied, err := hasAnyAppliedMigration(ctx, db)
	if err != nil || hasApplied {
		return err
	}
	exists, err := tableExists(ctx, db, "branches")
	if err != nil || !exists {
		return err
	}
	for _, file := range files {
		if strings.HasPrefix(file, "001_") {
			_, err := db.Exec(ctx, `INSERT INTO schema_migrations (version) VALUES ($1) ON CONFLICT DO NOTHING`, file)
			if err == nil {
				log.Printf("migration baselined for existing database: %s", file)
			}
			return err
		}
	}
	return nil
}

func hasAnyAppliedMigration(ctx context.Context, db *pgxpool.Pool) (bool, error) {
	var exists bool
	err := db.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM schema_migrations)`).Scan(&exists)
	return exists, err
}

func tableExists(ctx context.Context, db *pgxpool.Pool, table string) (bool, error) {
	var exists bool
	err := db.QueryRow(ctx, `SELECT to_regclass($1) IS NOT NULL`, "public."+table).Scan(&exists)
	return exists, err
}

func migrationApplied(ctx context.Context, db *pgxpool.Pool, version string) (bool, error) {
	var exists bool
	err := db.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM schema_migrations WHERE version=$1)`, version).Scan(&exists)
	return exists, err
}

func runMigration(ctx context.Context, db *pgxpool.Pool, version, sql string) error {
	tx, err := db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, sql); err != nil {
		return fmt.Errorf("migration %s failed: %w", version, err)
	}
	if _, err := tx.Exec(ctx, `INSERT INTO schema_migrations (version) VALUES ($1)`, version); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
