# GORMigro
GORMigro is tool for handling the gorm migrations with builtin cobra command for manual handle of migration process.

[![GoDoc](https://godoc.org/github.com/praggolabs/gormigro?status.svg)](https://godoc.org/github.com/praggolabs/gormigro)
[![GoReport](https://goreportcard.com/badge/praggolabs/gormigro)](https://goreportcard.com/report/praggolabs/gormigro)
[![Version](https://img.shields.io/badge/version-0.0.1-blue.svg)](https://github.com/praggolabs/grupttor/releases/latest)


# Table of Contents
- [Installing](#installing)
- [Configuration](#configuration)
  * [Schema initial function](#schema-init)
  * [Options](#options)
  * [Migration registration](#migration-registration)
- [Examples](#getting-started)
  * [Simple example](#simple-example)
- [Cobra command](#cobra-command)
- [Contributing](#contributing)
- [License](#license)

# Installing
Just run:
    `go get -u github.com/PragGoLabs/gormigro`

# Configuration

## Schema initial function
You can specify initial schema function which will handle the base of your schema, when you have an existing project
and you want to start with migrations and dont want to create single migration.

The function will be run on start, that means it will run when there's no executed migration - no records in migration table.
Also you've to enable `RunInitSchema` in options, which is in `gormirgo.DefaultOptions`.

Example of function:
```go
	// you can register initial schema function		
	g.RegisterInitialSchemaFunction(func(db *gorm.DB) error {

		type Car struct {
			Name string
			Color string
		}
		if db.HasTable(&Car{}) {
			db.DropTable(&Car{})
		}

		db.CreateTable(&Car{})

		return nil
	})
```

## Options
Gormigro have a 4 simple options
```go
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
```

You don't have to configure, just use `gormigro.DefaultOptions`.

## Migration registration
`Use migrations/ package as the place where you store all of your migrations`

then you can simply load all the migrations by:
```go
    import _ "{yourPackage}/migrations"
``` 

it will import all your migrations and run `init func()`.
Gormigro have a `gormigro.DefaultCollector` initialized in gormigro import by default,
so it collects all migrations.
  
There is 2 ways to register migration to the gormigro.

1. Use a init function in single migration file, example:
    ```go
        package migrations
        
        import (
            "errors"
            "github.com/PragGoLabs/gormigro"
            "github.com/jinzhu/gorm"
        )
        
        // and that's it, it will run on start and register to default collector
        func init() {
            type Items struct {
                gorm.Model
                
                Desc string
                Title string
            }
        
            gormigro.DefaultCollector.RegisterMigration(
                gormigro.NewMigration(
                    "Migration2019010301",
                    func(db *gorm.DB) error {
                        return db.CreateTable(&Items{}).Error
                    },
                    func(db *gorm.DB) error {
                        return db.DropTable(&Items{}).Error
                    },
                ),
            )
        }
    ``` 
2. if you need more logic, or you like to have migration in own struct, you can use the `Migrator` interface:
    ```go
        package migrations
        
        import (
            "github.com/PragGoLabs/gormigro"
            "github.com/jinzhu/gorm"
        )
        
        type Migration2019010301 struct {}
        
        func (m Migration2019010301) Identify() string {
            return "Migration2019010301"
        }
        
        func (m Migration2019010301) Migrate(db *gorm.DB) error {
            type Brands struct {
                gorm.Model
        
                Title string
                Uri string
            }
        
            return db.CreateTable(&Brands{}).Error
        }
        
        func (m Migration2019010301) Rollback(db *gorm.DB) error {
            type Brands struct {
                gorm.Model
        
                Title string
                Uri string
            }
        
            return db.DropTable(&Brands{}).Error
        }
        
        func init() {
            gormigro.DefaultCollector.RegisterMigration(Migration2019010301{})
        }

    ```
    collector will wrap these migration with `Migration` struct, and add to collection. 


# Examples

## Simple example with init function and inline migration create
```go
	// open your connection
	db, err := gorm.Open("mysql", "xxx")

	// configure gormigro with default options
	g := gormigro.NewGormigro(db, gormigro.DefaultOptions)

	// you can register initial schema function		
	g.RegisterInitialSchemaFunction(func(db *gorm.DB) error {

		type Car struct {
			Name string
			Color string
		}
		if db.HasTable(&Car{}) {
			db.DropTable(&Car{})
		}

		db.CreateTable(&Car{})

		return nil
	})
	
	// and inline register your migration
    g.AddMigration(gormigro.NewMigration(
        "Migration2019010301",
        func(db *gorm.DB) error {
            // do some logic
            return nil
        },
        func(db *gorm.DB) error {
            // do some logic
            return nil
        },
    ))
```
`We prefer and recommend you to use init register mechanism described in:` [Migration registration](#migration-registration)

# Cobra command
Gormigro have a built-in cobra command with migration interface.
There are 3 commands to use:
1. `migrate run`
    - it start the migration process
    - run all non executed migrations
2. `migrate clear`
    - rollback all executed migrations
    - and clear the migration table
3. `migrate drop`
    - it'll drop whole database schema, all tables which will be found
    - also with migration table
    
### Example of usage
```go
    // import the command
    import "github.com/PragGoLabs/gormigro/cmd"

    // create instance of gormigro
	g := gormigro.NewGormigro(db, gormigro.DefaultOptions)

    // and just register to your root cobra command with your gormigro instance 
    cmd.InitMigrationCommand(g, rootCmd)
```  
And that's it now you will be able your migration thru cmd interface, expected output:
```
    Usage:
      app migrate [CMD] [flags]
      app migrate [command]
    
    Available Commands:
      clear       Truncate migration table and rollback migrations, except initial
      drop        Drop all tables in schema also with migration table
      run         Run migrations
    
    Flags:
      -h, --help   help for migrate
    
    Use "app migrate [command] --help" for more information about a command.
```

# Contributing

1. Fork it
2. Clone (`git clone https://github.com/your_username/gormigro && cd gormigro`)
3. Create your feature branch (`git checkout -b my-new-feature`)
4. Do changes, Commit, Push
5. Create new pull request
6. Thanks in advance :-) 

# License

See [LICENSE.txt](https://github.com/praggolabs/gormigro/LICENSE.md)
