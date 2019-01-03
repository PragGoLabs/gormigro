package gormigro

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"log"
)

// Gormigro is migration tool which allow you to simply handle migrations with gorm.io
type Gormigro struct {
	// db gorm database
	db *gorm.DB

	// options refers to setup of gormigro
	options Options

	// initialSchemaMigration runs on start, when there's no executed migration
	initialSchemaMigration migrationFnc

	// migrationsCollection list of migrations
	migrationsCollection *Collection

	// migrationManager manage migration table
	migrationManager MigrationManager
}

// NewGormigro return instance of gormigro based on options passed
// work with all collected migrations
func NewGormigro(db *gorm.DB, options Options) *Gormigro {
	if options.DebugMode {
		db.LogMode(true)
	}

	return &Gormigro{
		db:                   db,
		options:              options,
		migrationsCollection: DefaultCollector.Export(),
		migrationManager:     NewMigrationManager(db.New(), options.MigrationTable),
	}
}

// NewGormigro return instance of gormigro based on options passed
// with specified migrations on input
func NewGormigroWithMigrations(db *gorm.DB, options Options, migrations []Migration) *Gormigro {
	if options.DebugMode {
		db.LogMode(true)
	}

	return &Gormigro{
		db:                   db,
		options:              options,
		migrationsCollection: NewCollectionWithMigrations(migrations),
		migrationManager:     NewMigrationManager(db.New(), options.MigrationTable),
	}
}

// AddMigration append migration to list of migrations
func (g *Gormigro) AddMigration(migration Migration) error {
	if g.migrationsCollection.Contains(migration) {
		return CreateMigrationAlreadyExistsError(migration.ID)
	}

	g.migrationsCollection.Append(migration)

	return nil
}

// RegisterInitialSchemaFunction attach init function to migration process
// It will run on start of migration process if there are no previously executed migration
func (g *Gormigro) RegisterInitialSchemaFunction(initFunc migrationFnc) {
	g.initialSchemaMigration = initFunc
}

// Migrate run the migration process
// Retrieve last executed migration and run from it
// If there's no executed migration, it run everything
// If you specify the initial schema function, it'll also run it if no previous executed migration
func (g *Gormigro) Migrate() error {
	lastMigration := g.migrationManager.GetLastExecutedMigration()
	var err error
	if lastMigration == nil && g.options.RunInitSchema && g.initialSchemaMigration != nil {
		err = g.initialSchemaMigration(g.db)
	}

	if err != nil {
		return err
	}

	if lastMigration == nil {
		// start from first
		return g.MigrateFrom("")
	}

	return g.MigrateFrom(lastMigration.MigrationId)
}

// MigrateFrom run migrations from specified one(except the specified)
func (g *Gormigro) MigrateFrom(id string) error {
	mc := g.migrationsCollection

	// sort collection by migrationId
	if g.options.SortMigrationsByIDField {
		mc = mc.SortBy(MigrationId, Asc)
	}

	// if id is not empty, reduce the collection from(not included) migration id
	if id != "" {
		mc = g.migrationsCollection.SliceFrom(id)
	}

	if mc.Empty() {
		log.Print("No migrations to execute")
	}

	for _, m := range mc.List() {
		log.Printf("Running migration with ID %s\n", m.ID)
		// start the transaction per migration
		tx := g.db.New()
		if tx.Error != nil {
			return tx.Error
		}

		// run migration
		err := m.Migrate(tx)

		if err != nil {
			log.Printf("Error occured when migration %s run [%s]\n", m.ID, err)
			log.Printf("Rollbacking %s\n", m.ID)

			m.Rollback(tx)

			return err
		}

		// and insert record to table
		g.migrationManager.AddMigration(m.ID)
	}

	if !mc.Empty() {
		log.Printf("%d migrations executed\n", mc.Count())
	}

	return nil
}

// Clear rollback all migrations and remove record of execution from migration table
func (g *Gormigro) Clear() error {
	// you have to rollback in reverse order
	sorted := g.migrationsCollection.SortBy(MigrationId, Asc)

	lastMigration := g.migrationManager.GetLastExecutedMigration()
	if lastMigration != nil {
		sorted = g.migrationsCollection.SliceFrom(lastMigration.MigrationId)
	}

	reverseOrder := sorted.SortBy(MigrationId, Desc)
	for _, m := range reverseOrder.List() {
		// just rollback the
		err := m.Rollback(g.db)
		if err != nil {
			return CreateUnableToRollbackMigrationError(m.ID)
		}

		g.migrationManager.RemoveMigration(m.ID)
	}

	return nil
}

// DropSchema drop the whole database schema
func (g *Gormigro) DropSchema() error {
	type Table struct {
		Name string
	}

	rows, err := g.db.DB().Query("SHOW TABLES")
	if err != nil {
		return err
	}

	defer rows.Close()

	// disable FK check for safe remove
	g.db.Exec("SET FOREIGN_KEY_CHECKS=0;")
	for rows.Next() {
		var table Table

		rows.Scan(&table.Name)
		log.Printf("Dropping table: %s\n", table.Name)

		res := g.db.Exec(fmt.Sprintf("DROP TABLE %s", table.Name))

		if res.Error != nil {
			return res.Error
		}
	}

	// and enable again
	g.db.Exec("SET FOREIGN_KEY_CHECKS=1;")

	log.Println("Dropping completed")

	return nil
}

// internal usage
func (g *Gormigro) runInitialSchemaFunc() error {
	log.Println("Running initial schema migration")

	// start transaction
	tx := g.db.Begin()

	defer tx.Close()

	// run init func
	// rollback if any problem occured
	if err := g.initialSchemaMigration(tx); err != nil {
		tx.Rollback()

		return err
	}

	return tx.Commit().Error
}
