package cmd

import (
	"github.com/spf13/cobra"
	"github.com/f4814n/piston/server"
	log "github.com/sirupsen/logrus"
)

var Debug bool
var Trace bool
var VeryTrace bool

var serverCmd = &cobra.Command{
	Use:	"server",
	Short:	"Run the server",
	Long:	"test",
	Run:	func(cmd *cobra.Command, args []string){
		if Trace {
			if VeryTrace {
				log.SetReportCaller(true)
			}
			log.SetLevel(log.TraceLevel)
		} else if Debug {
			log.SetLevel(log.DebugLevel)
		}

		s := server.Server{}
		s.Serve()
	},
}	

func init() {
	rootCmd.PersistentFlags().BoolVarP(&Debug, "debug", "d", false, "debug log")
	rootCmd.PersistentFlags().BoolVarP(&Trace, "trace", "", false, "trace log")
	rootCmd.PersistentFlags().BoolVarP(&VeryTrace, "verytrace", "", false, "very trace log")

	rootCmd.AddCommand(serverCmd)
}
