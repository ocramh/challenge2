package cli

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

var (
	item       string
	cid        string
	serverPort int
)

func init() {
	rootCmd.AddCommand(clientCmd)

	clientCmd.AddCommand(addItemCmd)
	addItemCmd.Flags().StringVarP(&item, "item", "i", "", "the item content")
	addItemCmd.Flags().IntVarP(&serverPort, "port", "p", 9999, "the http port the client will bind to")
	addItemCmd.MarkFlagRequired("item")

	clientCmd.AddCommand(getItemCmd)
	getItemCmd.Flags().StringVarP(&cid, "cid", "c", "", "the item content cid")
	getItemCmd.Flags().IntVarP(&serverPort, "port", "p", 9999, "the http port the client will bind to")
	getItemCmd.MarkFlagRequired("cid")
}

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "The client interface for interacting with a provider",
	Run:   func(cmd *cobra.Command, args []string) {},
}

var addItemCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds an item to the provider",
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := http.Get(fmt.Sprintf("http://localhost:%d/add?item=%s", serverPort, item))
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}

		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}

		fmt.Println(string(b))
	},
}

var getItemCmd = &cobra.Command{
	Use:   "get",
	Short: "Gets an item from the provider",
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := http.Get(fmt.Sprintf("http://localhost:%d/get?cid=%s", serverPort, cid))
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}

		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(b))
	},
}
