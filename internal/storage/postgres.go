package storage

import (
	"database/sql"
	"fmt"
	"website-dm/internal/config"
	"time"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func Open(cfg config.DBConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.DBName, cfg.Password, cfg.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("не вдалося відкрити БД: %w", err)
	}

	var lastErr error
	for i := 0; i < 12; i++ {
		if err = db.Ping(); err == nil {
			lastErr = nil
			break
		}
		lastErr = err
		time.Sleep(time.Duration(i+1) * 500 * time.Millisecond)
	}
	if lastErr != nil {
		return nil, fmt.Errorf("не вдалося підключитись до БД: %w", lastErr)
	}

	if err = goose.Up(db, "db/migrations"); err != nil {
		return nil, fmt.Errorf("помилка міграцій: %w", err)
	}

	return db, nil
}
