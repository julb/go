package cmd

import (
	"github.com/spf13/cobra"

	log "github.com/julb/go/pkg/logging"
)

func init() {
	// Command opts
	rootCmd.PersistentFlags().BoolVarP(&rootCmdOptInfo, "info", "", false, "Enables the INFO log level.")
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
			log.SetLevel("trace")
		} else if rootCmdOptDebug {
			log.SetLevel("debug")
		} else if rootCmdOptInfo {
			log.SetLevel("info")
		} else {
			log.SetLevel("warn")
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		// NOOP
	},
}

func ExecuteMainCmd() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatalf("Error when executing command %s", err)
	}
}
