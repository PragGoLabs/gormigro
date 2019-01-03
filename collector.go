package gormigro

import (
	"errors"
	"github.com/jinzhu/gorm"
)

// DefaultCollector is default instance live globally in context
var DefaultCollector = &Collector{
	defaultCollection: NewCollection(),
}

// Collector handle all existing migration
type Collector struct {
	// internal migrations handled
	defaultCollection *Collection
}

// RegisterMigration allow you to append migration to collector
func (c *Collector) RegisterMigration(v interface{}) error {
	var migration Migration

	switch m := v.(type) {
	case Migrator:
		migration = NewMigration(
			m.Identify(),
			func(db *gorm.DB) error {
				return m.Migrate(db)
			},
			func(db *gorm.DB) error {
				return m.Rollback(db)
			},
		)
		break
	case Migration:
		migration = m
		break
	default:
		return errors.New("unsupported migration type passed")
	}

	if c.defaultCollection.Contains(migration) {
		return CreateMigrationAlreadyExistsError(migration.ID)
	}

	c.defaultCollection.Append(migration)

	return nil
}

// Export collection with migrations
func (c *Collector) Export() *Collection {
	return c.defaultCollection
}
