package gormigro

import "gorm.io/gorm"

type migrationFnc func(*gorm.DB) error
type rollbackFnc func(*gorm.DB) error

// Migrator is interface that will allow you define your own migration struct
type Migrator interface {
	// Identify have to return ID of migration
	Identify() string

	// Migrate run migration body
	Migrate(*gorm.DB) error

	// Rollback run rollback body
	Rollback(*gorm.DB) error
}

// Migration is struct which represent single migration
type Migration struct {
	// ID is unique string which identify the migration
	ID string

	// Migrate run migration function
	Migrate migrationFnc

	// Rollback run rollback function
	Rollback rollbackFnc
}

// NewMigration return new migration instance
// its simple wrapper of two function and string identifier
func NewMigration(id string, migrate migrationFnc, rollback rollbackFnc) Migration {
	return Migration{
		ID:       id,
		Migrate:  migrate,
		Rollback: rollback,
	}
}
