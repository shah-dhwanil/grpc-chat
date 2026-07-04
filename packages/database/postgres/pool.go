package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	pgxUUID "github.com/vgarvardt/pgx-google-uuid/v5"
)

type PgPool = *pgxpool.Pool
type Tx = pgx.Tx

func NewPgPool(config *pgxpool.Config, pingTimeout *time.Duration) (PgPool,error) {
	afterConnect := config.AfterConnect
	config.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		// Set the timezone to UTC for the connection
		pgxUUID.Register(conn.TypeMap())
		if afterConnect != nil {
			return afterConnect(ctx, conn)
		}
		return nil
	}
	pgpool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}
	if pingTimeout == nil{
		return pgpool, nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), *pingTimeout)
	defer cancel()
	if err = pgpool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("Timeout Occured: Unable to connect to Postgres Database %v: %w",fmt.Sprintf("%s:%d/%s",config.ConnConfig.Host,config.ConnConfig.Port,config.ConnConfig.Database), err)
	}
	return pgpool, nil
}