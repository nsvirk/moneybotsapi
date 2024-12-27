package repository

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"

	"github.com/nsvirk/moneybotsapi/internal/config"
	"github.com/nsvirk/moneybotsapi/internal/models"
)

func ConnectPostgres(cfg *config.Config) (*bun.DB, error) {
	// Parse the DSN to validate and potentially modify it
	_, err := url.Parse(cfg.PostgresDsn)
	if err != nil {
		return nil, fmt.Errorf("invalid DSN: %v", err)
	}

	// Create connector with more robust settings
	connector := pgdriver.NewConnector(
		pgdriver.WithDSN(cfg.PostgresDsn),
		pgdriver.WithTimeout(5),
		pgdriver.WithDialTimeout(5*time.Second),
		pgdriver.WithReadTimeout(30*time.Second),
		pgdriver.WithWriteTimeout(30*time.Second),
		pgdriver.WithApplicationName(cfg.APIName),
		pgdriver.WithTLSConfig(nil),
	)

	// Open database with connection pool settings
	sqldb := sql.OpenDB(connector)
	sqldb.SetMaxOpenConns(4)
	sqldb.SetMaxIdleConns(2)
	sqldb.SetConnMaxLifetime(time.Hour)

	// Create Bun db instance
	db := bun.NewDB(sqldb, pgdialect.New())

	// Add custom query hook for logging
	db.AddQueryHook(NewDBLogger(cfg.PostgresLogLevel))

	// Test connection with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	log.Info().
		Str("module", "database").
		Msg("Connected to PostgreSQL")

	// Create schema if it doesn't exist
	if err := createSchema(ctx, db, cfg.PostgresSchema); err != nil {
		return nil, err
	}

	// Run migrations
	if err := runMigrations(ctx, db, cfg); err != nil {
		return nil, err
	}

	return db, nil
}

func createSchema(ctx context.Context, db *bun.DB, schema string) error {
	// Create schema with proper error handling
	query := fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", schema)
	if _, err := db.ExecContext(ctx, query); err != nil {
		return fmt.Errorf("failed to create schema: %v", err)
	}

	// Set search_path
	query = fmt.Sprintf("SET search_path TO %s,public", schema)
	if _, err := db.ExecContext(ctx, query); err != nil {
		return fmt.Errorf("failed to set search_path: %v", err)
	}

	return nil
}

func runMigrations(ctx context.Context, db *bun.DB, cfg *config.Config) error {
	tableName := fmt.Sprintf("%s.%s", cfg.PostgresSchema, models.AuthTableName)

	// Create auth table with explicit schema
	_, err := db.NewCreateTable().
		Model((*models.AuthModel)(nil)).
		IfNotExists().
		Table(tableName).
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("failed to create auth table %s: %v", tableName, err)
	}

	return nil
}
