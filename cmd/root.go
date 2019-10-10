package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	root     string
	logLevel uint
)
var rootCmd = &cobra.Command{
	Use:   "piston",
	Short: "Piston is a performant minecraft server",
	Long:  "Piston is a performant minecraft server",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func setLogLevel() {
	if logLevel < 7 {
		log.SetLevel(log.AllLevels[logLevel])
	} else if logLevel == 7 {
		log.SetLevel(log.TraceLevel)
		log.SetReportCaller(true)
	} else {
		log.WithFields(log.Fields{"level": logLevel}).Fatal("Unknown log level")
	}
}

func init() {
	cobra.OnInitialize(setLogLevel)

	rootCmd.PersistentFlags().StringVar(&root, "root", ".", "Server root")
	rootCmd.PersistentFlags().UintVar(&logLevel, "log", 4, "log level. 0 Is silent")
}
