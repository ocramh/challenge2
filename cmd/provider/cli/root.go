package cli

import (
	"log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "provider",
	Short: "",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf(err.Error())
	}
}
