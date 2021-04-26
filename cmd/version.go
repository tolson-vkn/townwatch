package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tolson-vkn/townwatch/common/version"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "print version",
	Long:  `print version`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Townwatch Version Info:")
		fmt.Printf("Version: %s\n", version.VersionString())
		fmt.Printf("Commit:  %s\n", version.GitCommitString())
	},
}
