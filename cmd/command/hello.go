package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

var name string

func init() {
	rootCmd.AddCommand(helloCmd)
	helloCmd.Flags().StringVarP(&name, "name", "n", "", "Name to say hello to")
}

var helloCmd = &cobra.Command{
	Use:   "hello",
	Short: "Print a hello message",
	Long:  `A longer description of the hello command.`,
	Run: func(cmd *cobra.Command, args []string) {
		if name != "" {
			fmt.Printf("Hello, %s!\n", name)
		} else {
			fmt.Println("Hello from Cobra!")
		}
	},
}
