package postgres

import (
	"embed"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shah-dhwanil/grpc-chat/apps/user-service/internal/config"
	"github.com/shah-dhwanil/grpc-chat/packages/database/postgres"
)

//go:embed migration/*.sql
var Migrations embed.FS

func NewPool(config config.Postgres) (postgres.PgPool,error) {
	poolConfig,err := pgxpool.ParseConfig(fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode= %s application_name=user-service", config.Host, config.Port, config.User, config.Password, config.Database, config.SSLMode))
	if err != nil {
		return nil, err
	}
	pool,err:= postgres.NewPgPool(poolConfig,nil)
	if err != nil {
		return nil, err
	}
	return pool,nil
}