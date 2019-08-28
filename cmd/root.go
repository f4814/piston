package cmd

import (
	"github.com/spf13/cobra"
	log "github.com/sirupsen/logrus"
)

var root string

var rootCmd = &cobra.Command{
	Use:	"piston",
	Short:	"Piston is a performant minecraft server",
	Long:	"Piston is a performant minecraft server",
	Run:	func(cmd *cobra.Command, args []string){
		log.Info("Not implemented")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Error(err)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&root, "root", ".", "Server root")
}
