package command

import (
	"fmt"
	"log"
	"net/http"

	"github.com/alirezaghasemi/user-manager/internal/config"
	"github.com/alirezaghasemi/user-manager/internal/container"
	"github.com/alirezaghasemi/user-manager/internal/delivary/http/handler"
	"github.com/alirezaghasemi/user-manager/internal/delivary/http/router"
	"github.com/alirezaghasemi/user-manager/internal/repository"
	"github.com/alirezaghasemi/user-manager/internal/usecase"
	"github.com/spf13/cobra"
)

// httpCmd represents the http command
var httpCmd = &cobra.Command{
	Use:   "http",
	Short: "Start http server",
	Run: func(cmd *cobra.Command, args []string) {
		initializeConfigs()

		log.Println("starting kyc http server")

		startServer(&Cfg)
	},
}

func init() {
	rootCmd.AddCommand(httpCmd)
}

func startServer(cfg *config.Config) error {
	c := container.NewContainer(*cfg)
	// ----- Repositories -----
	userRepository := repository.NewUserRepository(c.DB)

	// ----- Usecases -----
	userUsecase := usecase.NewUserUsecase(userRepository, c.Validate)

	// ----- Handlers -----
	userHandler := handler.NewUserHandler(userUsecase, c.Validate)

	// ----- Routers -----
	router := router.NewRouter(*userHandler)

	fmt.Printf("%s:%d\n", cfg.Server.Host, cfg.Server.Port)
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		WriteTimeout: cfg.Server.WriteTimeout,
		ReadTimeout:  cfg.Server.ReadTimeout,
		Handler:      router,
	}

	err := server.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}
