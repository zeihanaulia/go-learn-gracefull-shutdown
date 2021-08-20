package cmd

import (
	"github.com/spf13/cobra"
	"github.com/zeihanaulia/go-learn-gracefull-shutdown/handlers/rest"
)

// restCmd represents the rest command
var restCmd = &cobra.Command{
	Use:   "rest",
	Short: "A brief description of your command",
	RunE: func(cmd *cobra.Command, args []string) error {
		return rest.NewServer().Run()
	},
}

func init() {
	rootCmd.AddCommand(restCmd)
}
