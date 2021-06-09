package cmd

import (
	"fmt"

	"github.com/julb/go/pkg/build"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var (
	printVersionFormat = "%-25s%s\n"
)

func PrintVersion() error {
	fmt.Printf("%s %s %s\n", build.Info.Name, build.Info.Version, build.Info.Arch)

	fmt.Printf(printVersionFormat, "Built at:", build.Info.Time)
	fmt.Printf(printVersionFormat, "Group:", build.Info.Group)
	fmt.Printf(printVersionFormat, "Name:", build.Info.Name)
	fmt.Printf(printVersionFormat, "Artifact:", build.Info.Artifact)
	fmt.Printf(printVersionFormat, "Version:", build.Info.Version)
	fmt.Printf(printVersionFormat, "Arch:", build.Info.Arch)
	fmt.Printf(printVersionFormat, "Build version:", build.Info.BuildVersion)
	fmt.Printf(printVersionFormat, "Git short revision:", build.Info.GitShortRevision)
	fmt.Printf(printVersionFormat, "Git revision:", build.Info.GitRevision)

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
