package cmd

import (
	"fmt"

	"github.com/emgag/keyvault-certdeploy/internal/lib/version"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of keyvault-certdeploy",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("keyvault-certdeploy %s -- %s\n", version.Version, version.Commit)
	},
}
