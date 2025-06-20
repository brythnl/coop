package sqlc

import "github.com/jackc/pgx/v5/pgxpool"

type Store interface {
	Querier
}

type PostgresStore struct {
	*Queries
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) Store {
	return &PostgresStore{
		db:      db,
		Queries: New(db),
	}
}
