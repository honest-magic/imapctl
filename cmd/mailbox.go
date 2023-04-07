package cmd

import (
	"bitbucket.org/mis79/imapctl/utl"
	"fmt"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/spf13/cobra"
	"log"
)

var mailboxCmd = &cobra.Command{
	Use:   "mailbox",
	Short: "Mailbox related commands",
	Long:  "Mailbox related commands",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var mailboxListCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all available mailboxes",
	Long:  "List all available mailboxes",
	RunE:  listMailboxes,
}

func listMailboxes(cmd *cobra.Command, args []string) error {

	if verbose {
		log.Println("Connecting to server...")
	}

	// Connect to server
	addr := fmt.Sprintf("%s:%d", host, port)
	c, err := client.DialTLS(addr, nil)
	if err != nil {
		return err
	}

	if verbose {
		log.Println("Connected")
	}
	// Don't forget to logout
	defer c.Logout()

	// Login
	if err := c.Login(user, password); err != nil {
		return err
	}
	if verbose {
		log.Println("Logged in")
	}

	boxes := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)
	go func() {
		done <- c.List("", "*", boxes)
	}()

	var copy []*imap.MailboxInfo
	maxName := -1
	for box := range boxes {
		copy = append(copy, box)
		maxName = utl.Max(maxName, len(box.Name))
	}
	if err := <-done; err != nil {
		return err
	}

	fmtString := fmt.Sprintf("%%-%ds%%v\n", maxName+2)

	for _, box := range copy {
		fmt.Printf(fmtString, box.Name, box.Attributes)
	}

	return nil
}

func init() {
	mailboxCmd.AddCommand(mailboxListCmd)
	rootCmd.AddCommand(mailboxCmd)
}
