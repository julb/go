package cmd

import (
	"github.com/spf13/cobra"

	"github.com/julb/go/pkg/oas"
)

func init() {
	// Command opts
	oasIndexCmd.Flags().StringVarP(&oasIndexCmdOptDirectory, "directory", "d", ".", "Local directory containing the specifications to index.")
	oasIndexCmd.Flags().StringVarP(&oasIndexCmdOptUrl, "url", "u", "", "Public URL from which the local directory is reachable.")
	oasIndexCmd.Flags().StringArrayVarP(&oasIndexCmdOptExtensions, "extension", "e", []string{".json", ".yaml", ".yml"}, "File extensions of the specifications to consider when indexing.")

	// Build command hierarchy
	oasCmd.AddCommand(oasIndexCmd)
	rootCmd.AddCommand(oasCmd)
}

var oasIndexCmdOptDirectory string
var oasIndexCmdOptUrl string
var oasIndexCmdOptExtensions []string
var oasIndexCmd = &cobra.Command{
	Use:   "index",
	Short: "Index capabilities",
	Long:  `Index OAS3 specifications`,
	RunE: func(cmd *cobra.Command, args []string) error {
		options := oas.NewIndexOptions()
		options.Directory = oasIndexCmdOptDirectory
		options.Extensions = oasIndexCmdOptExtensions
		options.Url = oasIndexCmdOptUrl
		return oas.Index(options)
	},
}

var oasCmd = &cobra.Command{Use: "oas"}
