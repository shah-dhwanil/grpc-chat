package main

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/shah-dhwanil/grpc-chat/apps/user-service/internal/config"
	"github.com/shah-dhwanil/grpc-chat/apps/user-service/internal/postgres"
	pg "github.com/shah-dhwanil/grpc-chat/packages/database/postgres"
)


func main() {
	cfg:= config.GetConfig()
	connConfig,err := pgx.ParseConfig(fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode= %s application_name=user-service", cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.Database, cfg.Postgres.SSLMode))
	if err != nil {
		panic(err)
	}
	err= pg.Migrate(context.TODO(),connConfig,&postgres.Migrations)
	if err != nil {
		panic(err)
	}
}