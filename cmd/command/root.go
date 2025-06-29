package command

import (
	"os"

	"github.com/alirezaghasemi/user-manager/internal/config"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/cobra"
)

var (
	cfg     config.Config
	envFile string
	rootCmd = &cobra.Command{
		Use: "",
		Run: func(cmd *cobra.Command, args []string) {
			initializeConfigs()
		},
	}
)

func initializeConfigs() {
	if envFile != "" {
		err := godotenv.Load(envFile)
		if err != nil {
			panic(err)
		}
	} else {
		_ = godotenv.Load()
	}

	err := envconfig.Process("", &cfg)
	if err != nil {
		panic(err)
	}
}

func init() {
	cobra.OnInitialize()
	rootCmd.PersistentFlags().StringVarP(&envFile, "env-file", "e", ".env", ".env file")

	rootCmd.AddCommand(helloCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
