package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

func PrintVersion() error {
	fmt.Println("Version is X.X.X")
	return nil
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of j3oas",
	Long:  `All software has versions. This is J3OAS's`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := PrintVersion(); err != nil {
			return err
		}
		return nil
	},
}
