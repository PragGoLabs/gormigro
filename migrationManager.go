package gormigro

import (
	"gorm.io/gorm"
)

type MigrationManager struct {
	db        *gorm.DB
	tableName string
}

func NewMigrationManager(db *gorm.DB, tableName string) MigrationManager {
	mm := MigrationManager{
		tableName: tableName,
		db:        db,
	}

	if !mm.checkIfTableExists() {
		mm.initTable()
	}

	return mm
}

func (mm MigrationManager) AddMigration(id string) {
	record := NewMigrationTable(mm.tableName)
	record.MigrationId = id
	if mm.db.Error != nil {
		panic(mm.db.Error)
	}

	res := mm.db.Create(record)
	if res.Error != nil {
		panic(res.Error)
	}
}

func (mm MigrationManager) GetLastExecutedMigration() *MigrationTable {
	lastResult := mm.db.Last(NewMigrationTable(mm.tableName))
	if lastResult.RowsAffected == 0 {
		return nil
	}

	var lastMigration MigrationTable
	lastResult.Row().Scan(&lastMigration.MigrationId, &lastMigration.ExecutedAt)

	return &lastMigration
}

func (mm MigrationManager) RemoveMigration(id string) {
	record := NewMigrationTable(mm.tableName)
	record.MigrationId = id
	if mm.db.Error != nil {
		panic(mm.db.Error)
	}

	res := mm.db.Delete(record)
	if res.Error != nil {
		panic(res.Error)
	}
}

func (mm MigrationManager) ClearExecutedMigrations() error {
	res := mm.db.Exec("TRUNCATE `", mm.tableName)

	return res.Error
}

// initTable initialize migration table
func (mm MigrationManager) initTable() {
	mm.db.Migrator().CreateTable(NewMigrationTable(mm.tableName))
}

// checkIfTableExists check if migration table exists
func (mm MigrationManager) checkIfTableExists() bool {
	return mm.db.Migrator().HasTable(NewMigrationTable(mm.tableName))
}
