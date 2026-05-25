package postgres

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type Health struct {
	conn *sqlx.DB
}

func NewHealth(conn *sqlx.DB) Health {
	return Health{conn: conn}
}

func (h Health) Ping(ctx context.Context) error {
	return h.conn.PingContext(ctx)
}
