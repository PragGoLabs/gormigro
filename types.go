package gormigro

import (
	"time"
)

const (
	DefaultMigrationTable = "_migrations"
)

// Options is for setup gormigro
type Options struct {
	// MigrationTable specify name of table where executed migrations are stored
	// default is "_migrations"
	MigrationTable string

	// RunInitSchema tells if the gormigro may run the initial func with schema initialization
	// if there's no executed migration and specified func
	// default is true
	RunInitSchema bool

	// DebugMode means it'll print out whole process, times, logs
	DebugMode bool

	// SortByIDField means you can override the exception order of migrations
	// default is true, so its sorted by ID on start
	SortMigrationsByIDField bool
}

// DefaultOptions default settings for gormigro
var DefaultOptions = Options{
	MigrationTable:          DefaultMigrationTable,
	RunInitSchema:           true,
	DebugMode:               false,
	SortMigrationsByIDField: true,
}

// MigrationTable is used for creating
type MigrationTable struct {
	tableName string `gorm:"-"`

	MigrationId string     `gorm:"primary_key"`
	ExecutedAt  *time.Time `sql:"DEFAULT:current_timestamp"`
}

func NewMigrationTable() *MigrationTable {
	// temporary fix, gorm cloning the struct thru reflection and there's
	// so it'll override the tableName settings
	return &MigrationTable{
		tableName: DefaultMigrationTable,
	}
}

func (mt MigrationTable) TableName() string {
	return DefaultMigrationTable
}
