package command

import (
	"log"
	"os"

	"github.com/alirezaghasemi/user-manager/internal/config"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/cobra"
)

var (
	Cfg     config.Config
	envFile string
	rootCmd = &cobra.Command{
		Use: "user-manager",
		Run: func(cmd *cobra.Command, args []string) {
			initializeConfigs()
		},
	}
)

func initializeConfigs() {
	if envFile != "" {
		err := godotenv.Load(envFile)
		if err != nil {
			log.Fatalf("Error loading env file: %v", err)
		}
	} else {
		_ = godotenv.Load()
	}

	err := envconfig.Process("", &Cfg)
	if err != nil {
		log.Fatalf("Error parsing config: %v", err)
	}
}

func init() {
	cobra.OnInitialize(initializeConfigs)

	rootCmd.PersistentFlags().StringVarP(&envFile, "env-file", "e", ".env", ".env file")

	rootCmd.AddCommand(helloCmd)

	rootCmd.AddCommand(migrateCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
