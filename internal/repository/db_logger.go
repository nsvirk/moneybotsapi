package repository

import (
	"context"
	"regexp"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
)

type LogLevel string

const (
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelError LogLevel = "error"
)

// DBLogger is a custom query hook for logging database operations
type DBLogger struct {
	logLevel LogLevel
}

func NewDBLogger(level string) *DBLogger {
	return &DBLogger{
		logLevel: LogLevel(strings.ToLower(level)),
	}
}

func (h *DBLogger) BeforeQuery(ctx context.Context, event *bun.QueryEvent) context.Context {
	return ctx
}

func (h *DBLogger) AfterQuery(ctx context.Context, event *bun.QueryEvent) {
	// Skip logging for log table operations to prevent recursive logging
	if strings.Contains(event.Query, "\"logs\"") {
		return
	}

	duration := float64(time.Since(event.StartTime).Milliseconds())
	operation := getOperationType(event.Query)
	cleanQuery := cleanQueryString(event.Query)

	// Always log errors regardless of log level
	if event.Err != nil {
		log.Error().
			Str("module", "database").
			Str("operation", operation).
			Float64("duration_ms", duration).
			Str("query", cleanQuery).
			Err(event.Err).
			Msg("Database query error")
		return
	}

	// Log based on configured level
	switch h.logLevel {
	case LogLevelDebug:
		log.Debug().
			Str("module", "database").
			Str("operation", operation).
			Float64("duration_ms", duration).
			Str("query", cleanQuery).
			Msg("Database query executed")
	case LogLevelInfo:
		if isSchemaOperation(operation) || duration > 100 {
			log.Info().
				Str("module", "database").
				Str("operation", operation).
				Float64("duration_ms", duration).
				Msg("Database query completed")
		}
	}
}

// Helper function to get operation type
func getOperationType(query string) string {
	query = strings.TrimSpace(strings.ToUpper(query))
	for _, op := range []string{"SELECT", "INSERT", "UPDATE", "DELETE", "CREATE", "ALTER", "DROP", "SET"} {
		if strings.HasPrefix(query, op) {
			return op
		}
	}
	return "UNKNOWN"
}

// Helper function to identify schema operations
func isSchemaOperation(operation string) bool {
	schemaOps := map[string]bool{
		"CREATE": true,
		"ALTER":  true,
		"DROP":   true,
	}
	return schemaOps[operation]
}

// cleanQueryString cleans up SQL queries for logging
func cleanQueryString(query string) string {
	// Define regex pattern for escaped quotes in column names
	re := regexp.MustCompile(`\"([^\"]+)\"`)

	// Replace escaped quotes in column names with single quotes
	query = re.ReplaceAllString(query, "$1")

	// Remove any duplicate spaces
	query = strings.Join(strings.Fields(query), " ")

	return query
}
