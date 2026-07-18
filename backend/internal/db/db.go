package db

import (
	"context"

	"github.com/jackc/pgx/v5"
)


func Connect(ctx context.Context, url string) (*pgx.Conn, error) {
    return pgx.Connect(ctx, url)
}
