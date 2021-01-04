package gormigro

import (
	"errors"

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
	record := NewMigrationTable()
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
	lastResult := mm.db.Last(NewMigrationTable())
	if lastResult.RowsAffected == 0 {
		return nil
	}

	var lastMigration MigrationTable
	lastResult.Row().Scan(&lastMigration.MigrationId, &lastMigration.ExecutedAt)

	return &lastMigration
}

func (mm MigrationManager) RemoveMigration(id string) {
	record := NewMigrationTable()
	record.MigrationId = id
	if mm.db.Error != nil {
		panic(mm.db.Error)
	}

	res := mm.db.Delete(record)
	if res.Error != nil {
		panic(res.Error)
	}
}

func (mm MigrationManager) IsMigrationExecuted(id string) bool {
	res := mm.db.Find(NewMigrationTable())
	if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return false
	}
	if res.Error != nil {
		panic(res.Error)
	}

	return true
}

func (mm MigrationManager) ClearExecutedMigrations() error {
	res := mm.db.Exec("TRUNCATE `", mm.tableName)

	return res.Error
}

// initTable initialize migration table
func (mm MigrationManager) initTable() {
	mm.db.Migrator().CreateTable(NewMigrationTable())
}

// checkIfTableExists check if migration table exists
func (mm MigrationManager) checkIfTableExists() bool {
	return mm.db.Migrator().HasTable(NewMigrationTable())
}
