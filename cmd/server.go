package cmd

import (
	"github.com/f4814n/piston/server"
	"github.com/spf13/cobra"
)

var (
	address string
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run the server",
	Long:  "Run a piston server",
	Run: func(cmd *cobra.Command, args []string) {
		s := server.Server{}
		s.Serve(address)
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	rootCmd.PersistentFlags().StringVar(&address, "address", ":25565", "Server Address")
}
