package cmd

import (
	"fmt"
	"os"

	// homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	//"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:     "imapctl",
	Short:   "Control and manage IMAP based mailboxes",
	Long:    "Control and manage IMAP based mailboxes",
	Version: "0.0.1",
	RunE:    imapCtrl,
}

func init() {
	//rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
}

func imapCtrl(cmd *cobra.Command, args []string) error {
	return cmd.Help()
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
