package cmd

import (
	"fmt"

	"github.com/julb/go/pkg/build"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

func PrintVersion() error {
	fmt.Printf("%s %s %s\n", build.Info.Name, build.Info.Version, build.Info.Arch)
	fmt.Printf("Built at:			%s\n", build.Info.Time)
	fmt.Printf("Version:			%s\n", build.Version)
	fmt.Printf("Build version:		%s\n", build.Info.BuildVersion)
	fmt.Printf("Git revision:		%s\n", build.Info.GitRevision)
	fmt.Printf("Git short revision:	%s\n", build.Info.GitShortRevision)
	return nil
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information of the client",
	Long:  `Print build and version information about the j3 client`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := PrintVersion(); err != nil {
			return err
		}
		return nil
	},
}
