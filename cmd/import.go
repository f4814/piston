package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import Map",
	Long:  "test",
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Not implemented")
	},
}

func init() {
	rootCmd.AddCommand(importCmd)
}
