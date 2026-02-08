package postgres

import (
	"villainrsty-ecommerce-server/internal/adapters/persistence/postgres/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewQueries(pool *pgxpool.Pool) *sqlc.Queries {
	return sqlc.New(pool)
}
