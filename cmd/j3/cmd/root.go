package cmd

import (
	"github.com/spf13/cobra"

	log "github.com/sirupsen/logrus"
)

func init() {
	// Command opts
	rootCmd.PersistentFlags().BoolVarP(&rootCmdOptInfo, "info", "", false, "Enables the DEBUG log level.")
	rootCmd.PersistentFlags().BoolVarP(&rootCmdOptDebug, "debug", "", false, "Enables the DEBUG log level.")
	rootCmd.PersistentFlags().BoolVarP(&rootCmdOptTrace, "trace", "", false, "Enables the TRACE log level.")
}

var rootCmdOptInfo bool
var rootCmdOptDebug bool
var rootCmdOptTrace bool
var rootCmd = &cobra.Command{
	Use:   "j3",
	Short: "j3 is a command utility tool linked to julb.me projects",
	Long:  `blablabla`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if rootCmdOptTrace {
			log.SetLevel(log.TraceLevel)
		} else if rootCmdOptDebug {
			log.SetLevel(log.DebugLevel)
		} else if rootCmdOptInfo {
			log.SetLevel(log.InfoLevel)
		} else {
			log.SetLevel(log.WarnLevel)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		// NOOP
	},
}

func ExecuteMainCmd() error {
	return rootCmd.Execute()
}
