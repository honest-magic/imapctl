package cmd

import (
	"fmt"
	"log"

	"github.com/emersion/go-imap/client"
	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Checks connection to an IMAP mailbox",
	Long:  "Checks connection to an IMAP mailbox",
	RunE:  cmdCheck,
}

func cmdCheck(cmd *cobra.Command, args []string) error {

	// TODO: move settings to viper
	// Connect to server
	addr := fmt.Sprintf("%s:%d", host, port)
	c, err := client.DialTLS(addr, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Logout()
	return nil
}

func init() {
	rootCmd.AddCommand(checkCmd)
}
