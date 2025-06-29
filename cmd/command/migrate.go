package command

import (
	"fmt"
	"net/url"

	_ "github.com/amacneil/dbmate/pkg/driver/postgres"
	"github.com/amacneil/dbmate/v2/pkg/dbmate"
	"github.com/spf13/cobra"
	"gorm.io/driver/postgres"
)

func dbmateDB() *dbmate.DB {
	connStr := postgres.GetConnectionString(postgres.Config{
		Host:               cfg.Database.Host,
		Port:               cfg.Database.Port,
		Username:           cfg.Database.Username,
		Password:           cfg.Database.Password,
		Name:               cfg.Database.Name,
		MaxOpenConnections: cfg.Database.MaxOpenConnections,
	})

	u, _ := url.Parse(connStr)
	dbConn := dbmate.New(u)
	dbConn.FS = migrations.Migrations
	dbConn.MigrationsDir = []string{"./"}
	dbConn.AutoDumpSchema = false

	return dbConn
}

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("migrate called")
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// migrateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// migrateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
