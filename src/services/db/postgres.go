package db

import (
	"context"
	"database/sql"

	"github.com/harrybawsac/p1-go/src/models"
)

type PostgresAdapter struct {
	DB *sql.DB
}

func (p *PostgresAdapter) InsertReading(ctx context.Context, r models.Reading) error {
	// TODO: implement using upsert/unique constraints for idempotency
	_ = r
	return nil
}
