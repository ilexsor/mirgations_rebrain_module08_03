package migrations

import (
	"database/sql"
	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upAddUpdatedAt, downAddUpdatedAt)
}

func upAddUpdatedAt(tx *sql.Tx) error {
	_, err := tx.Exec(`
			ALTER TABLE users ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP DEFAULT now();

		`)
	if err != nil {
		return err
	}

	return nil
}

func downAddUpdatedAt(tx *sql.Tx) error {
	_, err := tx.Exec(`
			ALTER TABLE users DROP COLUMN IF EXISTS updated_at;
		`)
	if err != nil {
		return err
	}

	return nil
}
