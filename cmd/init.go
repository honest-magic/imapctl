package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes the imapctl settings",
	Long:  "Initializes the imapctl settings",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("init()")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
