package cmd

import (
	"github.com/PragGoLabs/gormigro"
	"github.com/spf13/cobra"
	"log"
)

var gormiGro *gormigro.Gormigro

var rootMigrateCmd = &cobra.Command{
	Use:   "migrate [CMD]",
	Short: "Run database migrations",
	Long:  `Run database migration`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
	},
}

var migrationRunCmd = &cobra.Command{
	Use: "run",
	Short: "Run migrations",
	Run: func(cmd *cobra.Command, args []string) {
		// run migrate
		err := gormiGro.Migrate()

		if err != nil {
			log.Panic(err)
		}
	},
}

var migrationClearSchemaCmd = &cobra.Command{
	Use: "clear",
	Short: "Truncate migration table and rollback migrations, except initial",
	Run: func(cmd *cobra.Command, args []string) {
		// run migrate
		err := gormiGro.Clear()

		if err != nil {
			log.Panic(err)
		}
	},
}

var migrationDropSchemaCmd = &cobra.Command{
	Use: "drop",
	Short: "Drop all tables in schema also with migration table",
	Run: func(cmd *cobra.Command, args []string) {
		err := gormiGro.DropSchema()

		if err != nil {
			log.Panic(err)
		}
	},
}

// InitMigrationCommand register predefined migration cmd with iface to your cobra root command
func InitMigrationCommand(gormi *gormigro.Gormigro, rootCmd *cobra.Command) {
	// assign the control struct
	gormiGro = gormi

	// register sub commands
	rootMigrateCmd.AddCommand(migrationRunCmd, migrationClearSchemaCmd, migrationDropSchemaCmd)

	// register to root cmd
	rootCmd.AddCommand(rootMigrateCmd)
}