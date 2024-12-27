package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

var (
	logOnce sync.Once
)

type APILog struct {
	ID        uint      `gorm:"primarykey"`
	Timestamp time.Time `gorm:"index"`
	Level     string
	Message   string
	Module    string `gorm:"index"`
	Data      string `gorm:"type:jsonb"`
}

func (APILog) TableName() string {
	return "logs"
}

type DBWriter struct {
	db *gorm.DB
}

func NewDBWriter(db *gorm.DB) *DBWriter {
	return &DBWriter{db: db}
}

func (w *DBWriter) Write(p []byte) (n int, err error) {
	var event map[string]interface{}
	if err := json.Unmarshal(p, &event); err != nil {
		return 0, err
	}

	// Extract standard fields
	now := time.Now()
	level := fmt.Sprintf("%v", event["level"])
	msg := fmt.Sprintf("%v", event["message"])
	module := ""
	if m, ok := event["module"].(string); ok {
		module = m
	}

	// Create log entry
	logEntry := APILog{
		Timestamp: now,
		Level:     level,
		Message:   msg,
		Module:    module,
		Data:      string(p),
	}

	// Write to database synchronously to ensure it works
	if err := w.db.Create(&logEntry).Error; err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write log to database: %v\n", err)
		return 0, err
	}

	return len(p), nil
}

func InitLogger(db *gorm.DB) error {
	// Auto migrate the logs table
	if err := db.AutoMigrate(&APILog{}); err != nil {
		return err
	}

	// Set up file logging
	logFile, err := getDailyLogFile()
	if err != nil {
		return fmt.Errorf("failed to create log file: %v", err)
	}

	// Configure zerolog - Format RFC3339
	zerolog.TimeFieldFormat = time.RFC3339

	// Create console writer with colors
	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
		NoColor:    false,
		FormatLevel: func(i interface{}) string {
			var l string
			if ll, ok := i.(string); ok {
				switch ll {
				case "debug":
					l = "\033[36m" // cyan
				case "info":
					l = "\033[34m" // blue
				case "warn":
					l = "\033[33m" // yellow
				case "error":
					l = "\033[31m" // red
				case "fatal", "panic":
					l = "\033[35m" // magenta
				default:
					l = "\033[37m" // white
				}
				return fmt.Sprintf("%s%-6s\033[0m", l, strings.ToUpper(ll))
			}
			return "????"
		},
		FormatMessage: func(i interface{}) string {
			return fmt.Sprintf("%s", i)
		},
		FormatFieldName: func(i interface{}) string {
			return fmt.Sprintf("\033[32m%s\033[0m=", i)
		},
		FormatFieldValue: func(i interface{}) string {
			return fmt.Sprintf("%s", i)
		},
	}

	// Create DB writer
	dbWriter := NewDBWriter(db)

	// Create the multi writer
	multiWriter := zerolog.MultiLevelWriter(consoleWriter, logFile, dbWriter)

	// Create and set the default logger
	logger := zerolog.New(multiWriter).With().Timestamp().Logger()

	// Set as default logger
	log.Logger = logger

	// Start log rotation in background
	go rotateLogDaily(db)

	return nil
}

func getLogDirectory() string {
	if os.Getenv("MB_API_ENV") == "production" {
		return "/var/log/moneybotsapi"
	}
	return "_logs"
}

func getDailyLogFile() (*os.File, error) {
	logDir := getLogDirectory()
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %v", err)
	}

	filename := fmt.Sprintf("moneybotsapi-%s.log", time.Now().Format("2006-01-02"))
	logPath := filepath.Join(logDir, filename)

	return os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
}

func rotateLogDaily(db *gorm.DB) {
	for {
		now := time.Now()
		next := now.Add(24 * time.Hour)
		next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())
		duration := next.Sub(now)

		time.Sleep(duration)

		logFile, err := getDailyLogFile()
		if err != nil {
			log.Error().
				Err(err).
				Msg("Failed to rotate log file")
			continue
		}

		// Create console writer with colors
		consoleWriter := zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
			NoColor:    false,
			FormatLevel: func(i interface{}) string {
				var l string
				if ll, ok := i.(string); ok {
					switch ll {
					case "debug":
						l = "\033[36m" // cyan
					case "info":
						l = "\033[34m" // blue
					case "warn":
						l = "\033[33m" // yellow
					case "error":
						l = "\033[31m" // red
					case "fatal", "panic":
						l = "\033[35m" // magenta
					default:
						l = "\033[37m" // white
					}
					return fmt.Sprintf("%s%-6s\033[0m", l, strings.ToUpper(ll))
				}
				return "????"
			},
			FormatMessage: func(i interface{}) string {
				return fmt.Sprintf("%s", i)
			},
			FormatFieldName: func(i interface{}) string {
				return fmt.Sprintf("\033[32m%s\033[0m=", i)
			},
			FormatFieldValue: func(i interface{}) string {
				return fmt.Sprintf("%s", i)
			},
		}

		// Create DB writer
		dbWriter := NewDBWriter(db)

		// Create the multi writer
		multiWriter := zerolog.MultiLevelWriter(consoleWriter, logFile, dbWriter)

		// Create new logger
		logger := zerolog.New(multiWriter).With().Timestamp().Logger()

		// Update the default logger
		log.Logger = logger
	}
}
