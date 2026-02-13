package sql

import (
	"context"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

type (
	Config struct {
		Username     string `json:"username"`
		Password     string `json:"password"`
		URL          string `json:"url"`
		Port         string `json:"port"`
		DatabaseName string `json:"database_name"`
	}
)

func New(ctx context.Context, cfg Config) (*sqlx.DB, error) {
	if cfg.Username == "" || cfg.Password == "" || cfg.DatabaseName == "" {
		err := errors.New("[pkg SQL] invalid params, param is empty")
		return nil, err
	}

	if cfg.Port == "" {
		cfg.Port = "5432"
	}

	urlConn := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable host=%s port=%s", cfg.Username, cfg.Password, cfg.DatabaseName, cfg.URL, cfg.Port)
	db, err := sqlx.Open("postgres", urlConn)
	if err != nil {
		log.Printf("[pkg SQL] Unable to connect to database %s:%s: %s\n", cfg.URL, cfg.DatabaseName, err.Error())
		return nil, err
	}
	return db, nil
}
