// Package note provides access to the note table in the MySQL database.
package note

import (
	"database/sql"
	"fmt"

	"github.com/blue-jay/blueprint/model"
	database "github.com/blue-jay/core/storage/driver/mysql"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var (
	// table is the table name.
	table = "note"
)

// Item defines the model.
type Item struct {
	ID        uint32         `db:"id"`
	Name      string         `db:"name"`
	UserID    uint32         `db:"user_id"`
	CreatedAt mysql.NullTime `db:"created_at"`
	UpdatedAt mysql.NullTime `db:"updated_at"`
	DeletedAt mysql.NullTime `db:"deleted_at"`
}

// Configuration defines the shared configuration interface.
type Configuration struct {
	DB *sqlx.DB
}

// Config returns the global connection information.
func Config() Configuration {
	return Configuration{
		DB: database.SQL,
	}
}

// ByID gets an item by ID.
func (c Configuration) ByID(ID string, userID string) (Item, error) {
	result := Item{}
	err := c.DB.Get(&result, fmt.Sprintf(`
		SELECT id, name, user_id, created_at, updated_at, deleted_at
		FROM %v
		WHERE id = ?
			AND user_id = ?
			AND deleted_at IS NULL
		LIMIT 1
		`, table),
		ID, userID)
	return result, model.StandardError(err)
}

// ByUserID gets all entities for a user.
func (c Configuration) ByUserID(userID string) ([]Item, error) {
	var result []Item
	err := c.DB.Select(&result, fmt.Sprintf(`
		SELECT id, name, user_id, created_at, updated_at, deleted_at
		FROM %v
		WHERE user_id = ?
			AND deleted_at IS NULL
		`, table),
		userID)
	return result, model.StandardError(err)
}

// Create adds an item.
func (c Configuration) Create(name string, userID string) (sql.Result, error) {
	result, err := c.DB.Exec(fmt.Sprintf(`
		INSERT INTO %v
		(name, user_id)
		VALUES
		(?,?)
		`, table),
		name, userID)
	return result, model.StandardError(err)
}

// Update makes changes to an existing item.
func (c Configuration) Update(name string, ID string, userID string) (sql.Result, error) {
	result, err := c.DB.Exec(fmt.Sprintf(`
		UPDATE %v
		SET name = ?
		WHERE id = ?
			AND user_id = ?
			AND deleted_at IS NULL
		LIMIT 1
		`, table),
		name, ID, userID)
	return result, model.StandardError(err)
}

// DeleteHard removes an item.
func (c Configuration) DeleteHard(ID string, userID string) (sql.Result, error) {
	result, err := c.DB.Exec(fmt.Sprintf(`
		DELETE FROM %v
		WHERE id = ?
			AND user_id = ?
			AND deleted_at IS NULL
		`, table),
		ID, userID)
	return result, model.StandardError(err)
}

// DeleteSoft marks an item as removed.
func (c Configuration) DeleteSoft(ID string, userID string) (sql.Result, error) {
	result, err := c.DB.Exec(fmt.Sprintf(`
		UPDATE %v
		SET deleted_at = NOW()
		WHERE id = ?
			AND user_id = ?
			AND deleted_at IS NULL
		LIMIT 1
		`, table),
		ID, userID)
	return result, model.StandardError(err)
}
