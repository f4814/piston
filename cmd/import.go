package cmd

import (
	"github.com/spf13/cobra"
	log "github.com/sirupsen/logrus"
)

var importCmd = &cobra.Command{
	Use:	"import",
	Short:	"Import Map",
	Long:	"test",
	Run:	func(cmd *cobra.Command, args []string){
		log.Info("Not implemented")
	},
}	

func init() {
	rootCmd.AddCommand(importCmd)
}
