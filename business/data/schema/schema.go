package schema

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/ardanlabs/darwin"
	"github.com/jmoiron/sqlx"
	"github.com/shalimski/ultimateservice/business/sys/database"
)

var (
	//go:embed sql/schema.sql
	schemaDoc string

	//go:embed sql/delete.sql
	deleteDoc string
)

// Migrate attempts to bring the schema for db up to date with the migrations
// defined in this package.
func Migrate(ctx context.Context, db *sqlx.DB) error {
	if err := database.StatusCheck(ctx, db); err != nil {
		return fmt.Errorf("status check database: %w", err)
	}

	driver, err := darwin.NewGenericDriver(db.DB, darwin.PostgresDialect{})
	if err != nil {
		return fmt.Errorf("construct darwin driver: %w", err)
	}

	d := darwin.New(driver, darwin.ParseMigrations(schemaDoc))
	return d.Migrate()
}
