package migrations

import (
	"database/sql"
	"fmt"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upSplitName, downSplitName)
}

func upSplitName(tx *sql.Tx) error {

	// Добавляем новые колонки
	_, err := tx.Exec(`
		ALTER TABLE users
		ADD COLUMN IF NOT EXISTS last_name VARCHAR(100);
	`)
	if err != nil {
		return fmt.Errorf("failed to add columns: %w", err)
	}

	_, err = tx.Exec(`
	UPDATE users 
		SET name = split_part(name, ' ', 1), last_name = split_part(name, ' ', 2), updated_at = NOW()
		`)

	if err != nil {
		return fmt.Errorf("failed to update user : %w", err)
	}

	return nil

}

func downSplitName(tx *sql.Tx) error {

	_, err := tx.Exec(`
			UPDATE users 
			SET name = CONCAT(name, ' ', last_name), updated_at = NOW()
		`)
	if err != nil {
		return fmt.Errorf("failed to update user %w", err)
	}

	// Удаляем новые колонки
	_, err = tx.Exec(`
		ALTER TABLE users
		DROP COLUMN IF EXISTS last_name;
	`)
	if err != nil {
		return fmt.Errorf("failed to drop columns: %w", err)
	}

	return nil
}
