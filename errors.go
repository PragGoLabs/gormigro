package gormigro

import (
	"fmt"
)

// DefaultError implements errors interface to consist error handling
type DefaultError struct {
	Message string
}

// implemented error interface
func (de *DefaultError) Error() string {
	return de.Message
}

// MigrationAlreadyExistsError error wrap struct
type MigrationAlreadyExistsError struct {
	*DefaultError

	// MigrationId id of wrong migration
	MigrationId string
}

// CreateMigrationAlreadyExistsError factory for error create
func CreateMigrationAlreadyExistsError(migrationId string) MigrationAlreadyExistsError {
	return MigrationAlreadyExistsError{
		DefaultError: &DefaultError{
			Message: fmt.Sprintf("Migration with id %s already exists", migrationId),
		},
		MigrationId: migrationId,
	}
}

// UnableToRollbackMigrationError error wrap struct
type UnableToRollbackMigrationError struct {
	*DefaultError

	// MigrationId id of wrong migration
	MigrationId string
}

func CreateUnableToRollbackMigrationError(migrationId string) UnableToRollbackMigrationError {
	return UnableToRollbackMigrationError{
		DefaultError: &DefaultError{
			Message: fmt.Sprintf("Unable to rollback migration with %s", migrationId),
		},
		MigrationId: migrationId,
	}
}