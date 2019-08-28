package cmd

import (
	"github.com/spf13/cobra"
	log "github.com/sirupsen/logrus"
)

var exportCmd = &cobra.Command{
	Use:	"export",
	Short:	"Export map",
	Long:	"export",
	Run:	func(cmd *cobra.Command, args []string){
		log.Info("Not implemented")
	},
}	

func init() {
	rootCmd.AddCommand(exportCmd)
}
