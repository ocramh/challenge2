package cli

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/ocramh/challenge2/pkg/provider"
	"github.com/ocramh/challenge2/pkg/server"
)

var (
	capacity int
	dataDir  string
	httpPort int
)

func init() {
	rootCmd.AddCommand(providerCmd)
	providerCmd.Flags().IntVarP(&capacity, "capacity", "c", 5, "storage capacity (number of items)")
	providerCmd.Flags().StringVarP(&dataDir, "datadir", "d", "./data", "the directory for storing content")
	providerCmd.Flags().IntVarP(&httpPort, "port", "p", 9999, "the http port the provider server will bind to")
}

var providerCmd = &cobra.Command{
	Use:   "server",
	Short: "Runs a content provider behind an HTTP server",
	Run: func(cmd *cobra.Command, args []string) {
		prv, err := provider.New(capacity, dataDir)
		if err != nil {
			log.Fatalf(err.Error())
		}

		svr := server.NewLocalServer(httpPort, prv)
		svr.Start()
	},
}
