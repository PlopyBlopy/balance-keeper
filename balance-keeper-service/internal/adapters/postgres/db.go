package postgres

import "github.com/jackc/pgx/v5/pgxpool"

type db struct {
	pool *pgxpool.Pool
}
