package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const version = "Wali robot v0.1 -- HEAD"

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Wali",
	Long:  `The version number of Wali`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version)
	},
}
