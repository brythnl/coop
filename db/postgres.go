package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectDB(url string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, err
	}
	
	pool, err := pgxpool.New(context.Background(), config.ConnString())
	if err != nil {
		return nil, err
	}

	if err = pool.Ping(context.Background()); err != nil {
		return nil, err
	}

	log.Println("successfully connected to the database")
	return pool, nil
}
