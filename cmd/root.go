package cmd

import (
	"fmt"
	"log"
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

var host string
var port int16

var tls bool

var user string
var password string

var verbose bool

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().StringVar(&host, "host", "", "IMAP mail host (required)")
	rootCmd.PersistentFlags().Int16VarP(&port, "port", "p", 993, "IMAP port")
	rootCmd.PersistentFlags().BoolVarP(&tls, "tls", "s", true, "IMAP dial tls")
	rootCmd.PersistentFlags().StringVarP(&user, "user", "u", "", "IMAP user name (required)")
	rootCmd.PersistentFlags().StringVarP(&password, "password", "x", "", "IMAP user password (required)")

	err := rootCmd.MarkPersistentFlagRequired("host")
	if err != nil {
		log.Fatal(err)
	}
	err = rootCmd.MarkPersistentFlagRequired("user")
	if err != nil {
		log.Fatal(err)
	}
	err = rootCmd.MarkPersistentFlagRequired("password")
	if err != nil {
		log.Fatal(err)
	}

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
