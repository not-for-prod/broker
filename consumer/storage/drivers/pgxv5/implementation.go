package pgxv5

import "github.com/jackc/pgx/v5/pgxpool"

type Implementation struct {
	pool *pgxpool.Pool
}

func NewImplementation(pool *pgxpool.Pool) *Implementation {
	return &Implementation{
		pool: pool,
	}
}
