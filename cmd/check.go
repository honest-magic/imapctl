package cmd

import (
	"log"

	imap "github.com/emersion/go-imap"
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
	// TODO: ensure existence of mailboxes INBOX.Archive.YEAR.Sent/Inbox
	// TODO: move all messages older than a certain time to the archive

	log.Println("Connecting to server...")

	// Connect to server
	c, err := client.DialTLS("paragon.sui-inter.net:993", nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected")

	// Don't forget to logout
	defer c.Logout()

	// Login
	if err := c.Login("michael@szediwy.ch", "###Katie1979"); err != nil {
		log.Fatal(err)
	}
	log.Println("Logged in")

	// List mailboxes
	mailboxes := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)
	go func() {
		done <- c.List("", "*", mailboxes)
	}()

	if err := c.Create("INBOX.Archiv.2023"); err != nil {
		log.Fatal(err)
	}

	log.Println("Mailboxes:")
	for m := range mailboxes {
		log.Println("* " + m.Name)
	}

	if err := <-done; err != nil {
		log.Fatal(err)
	}

	// Select INBOX
	mbox, err := c.Select("INBOX", false)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Flags for INBOX:", mbox.Flags)

	// Get the last 4 messages
	from := uint32(1)
	to := mbox.Messages
	if mbox.Messages > 3 {
		// We're using unsigned integers here, only subtract if the result is > 0
		from = mbox.Messages - 3
	}
	seqset := new(imap.SeqSet)
	seqset.AddRange(from, to)

	messages := make(chan *imap.Message, 10)
	done = make(chan error, 1)
	go func() {
		done <- c.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, messages)
	}()

	log.Println("Last 4 messages:")
	for msg := range messages {
		log.Println("* " + msg.Envelope.Subject)
	}

	if err := <-done; err != nil {
		log.Fatal(err)
	}

	log.Println("Done!")
	return nil
}

func init() {
	rootCmd.AddCommand(checkCmd)
}
