package postgres

import (
	"context"
	"embed"
	"fmt"
	"io/fs"


	"github.com/jackc/pgx/v5"
	tern "github.com/jackc/tern/v2/migrate"
)



func Migrate(ctx context.Context, config *pgx.ConnConfig, migrations *embed.FS) error {

	conn, err := pgx.ConnectConfig(ctx,config)
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	m, err := tern.NewMigrator(ctx, conn, "schema_version")
	if err != nil {
		return fmt.Errorf("constructing database migrator: %w", err)
	}
	subtree, err := fs.Sub(migrations, "migration")
	if err != nil {
		return fmt.Errorf("retrieving database migrations subtree: %w", err)
	}
	if err := m.LoadMigrations(subtree); err != nil {
		return fmt.Errorf("loading database migrations: %w", err)
	}
	from, err := m.GetCurrentVersion(ctx)
	if err != nil {
		return fmt.Errorf("retreiving current database migration version")
	}
	if err := m.Migrate(ctx); err != nil {
		return err
	}
	if from == int32(len(m.Migrations)) {
		fmt.Printf("Database uptodate. %d",from)
	} else {
		fmt.Printf("Database migrated successfully %d %d.",from,len(m.Migrations))
	}
	return nil
}