package command

import (
	"fmt"
	"log"
	"os"
	"strconv"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/pressly/goose/v3"
	"github.com/spf13/cobra"
)

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Manage database migrations",
	Long:  `This command allows you to run database migrations using Goose`,
}

// up command
var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Run all pending migrations",
	Run: func(cmd *cobra.Command, args []string) {
		runMigration("up", 0)
	},
}

// up command
var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Rollback the last migration",
	Run: func(cmd *cobra.Command, args []string) {
		runMigration("down", 0)
	},
}

// status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show the status of migrations",
	Run: func(cmd *cobra.Command, args []string) {
		runMigration("status", 0)
	},
}

var redoCmd = &cobra.Command{
	Use:   "redo",
	Short: "Rollback and re-apply the last migration",
	Run: func(cmd *cobra.Command, args []string) {
		runMigration("redo", 0)
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the current migration version",
	Run: func(cmd *cobra.Command, args []string) {
		runMigration("version", 0)
	},
}

var fixCmd = &cobra.Command{
	Use:   "fix",
	Short: "Rename migrations to match sequential ordering",
	Run: func(cmd *cobra.Command, args []string) {
		err := goose.Fix("migrations")
		if err != nil {
			log.Fatalf("goose fix failed: %v", err)
		}
		fmt.Println("Migration files fixed successfully")
	},
}

var createCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a new migration file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := goose.Create(nil, "migrations", args[0], "sql")
		if err != nil {
			log.Fatalf("goose create failed: %v", err)
		}
		fmt.Printf("Migration %s created successfully\n", args[0])
	},
}

var upToCmd = &cobra.Command{
	Use:   "up-to [version]",
	Short: "Migrate up to a specific version",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		version, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			log.Fatalf("invalid version: %v", err)
		}
		runMigration("up-to", version)
	},
}

var downToCmd = &cobra.Command{
	Use:   "down-to [version]",
	Short: "Rollback down to a specific version",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		version, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			log.Fatalf("invalid version: %v", err)
		}
		runMigration("down-to", version)
	},
}

func init() {
	migrateCmd.AddCommand(upCmd)
	migrateCmd.AddCommand(downCmd)
	migrateCmd.AddCommand(statusCmd)
	migrateCmd.AddCommand(redoCmd)
	migrateCmd.AddCommand(versionCmd)
	migrateCmd.AddCommand(fixCmd)
	migrateCmd.AddCommand(createCmd)
	migrateCmd.AddCommand(upToCmd)
	migrateCmd.AddCommand(downToCmd)
}

func runMigration(command string, version int64) {
	// directory sql files

	goose.SetBaseFS(os.DirFS("."))

	dbString := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable", Cfg.Database.Username, Cfg.Database.Password, Cfg.Database.Host, Cfg.Database.Port, Cfg.Database.Name)
	// dbString := "postgresql://postgres:admin@127.0.0.1:5432/user_manager?sslmode=disable"
	fmt.Println(dbString)
	db, err := goose.OpenDBWithDriver("postgres", dbString)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	// run goose statements
	switch command {
	case "up":
		goose.Up(db, "migrations")
		if err != nil {
			log.Fatalf("goose up failed: %v", err)
		}
		fmt.Println("Migration up completed successfully")
	case "down":
		err = goose.Down(db, "migrations")
		if err != nil {
			log.Fatalf("goose down failed: %v", err)
		}
		fmt.Println("Migration down completed successfully")
	case "status":
		err = goose.Status(db, "migrations")
		if err != nil {
			log.Fatalf("goose status failed: %v", err)
		}
	case "redo":
		err = goose.Redo(db, "migrations")
	case "version":
		var v int64
		v, err = goose.GetDBVersion(db)
		if err == nil {
			fmt.Printf("Current DB version: %d\n", v)
		}
	case "up-to":
		err = goose.UpTo(db, "migrations", version)
	case "down-to":
		err = goose.DownTo(db, "migrations", version)
	}

	if err != nil {
		log.Fatalf("goose %s failed: %v", command, err)
	}
}
