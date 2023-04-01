package cmd

import (
	"fmt"
	imap "github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/spf13/cobra"
	"log"
)

var archiveCmd = &cobra.Command{
	Use:   "archive",
	Short: "Archives mail messages into an archive folder",
	Long:  "Archives mail messages into an archive folder",
	RunE:  cmdArchive,
}

var host string
var port int16

var tls bool

var user string
var password string

var target string
var inbox string
var outbox []string

func init() {

	archiveCmd.Flags().StringVar(&host, "host", "", "IMAP mail host (required)")
	archiveCmd.Flags().Int16VarP(&port, "port", "p", 993, "IMAP port")
	archiveCmd.Flags().BoolVarP(&tls, "tls", "s", true, "IMAP dial tls")
	archiveCmd.Flags().StringVarP(&user, "user", "u", "", "IMAP user name (required)")
	archiveCmd.Flags().StringVarP(&password, "password", "x", "", "IMAP user password (required)")
	archiveCmd.Flags().StringVarP(&target, "target", "t", "Archive", "The target mailbox for the archive")
	archiveCmd.Flags().StringVarP(&inbox, "inbox", "i", "INBOX", "The mailbox name of received messages")
	archiveCmd.Flags().StringSliceVarP(&outbox, "outbox", "o", []string{"INBOX.Sent"}, "The mailbox name of sent messages")

	err := archiveCmd.MarkFlagRequired("host")
	if err != nil {
		log.Fatal(err)
	}
	err = archiveCmd.MarkFlagRequired("user")
	if err != nil {
		log.Fatal(err)
	}
	err = archiveCmd.MarkFlagRequired("password")
	if err != nil {
		log.Fatal(err)
	}

	rootCmd.AddCommand(archiveCmd)
}

func cmdArchive(cmd *cobra.Command, args []string) error {

	log.Println("Connecting to server...")

	// Connect to server
	addr := fmt.Sprintf("%s:%d", host, port)
	c, err := client.DialTLS(addr, nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected")

	// Don't forget to logout
	defer c.Logout()

	// Login
	if err := c.Login(user, password); err != nil {
		log.Fatal(err)
	}
	log.Println("Logged in")

	// List mailboxes
	mailboxes := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)
	go func() {
		done <- c.List("", "*", mailboxes)
	}()

	//if err := c.Create("INBOX.Archiv.2023"); err != nil {
	//	log.Fatal(err)
	//}

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
	//from := uint32(1)
	to := mbox.Messages
	//if mbox.Messages > 3 {
	//	// We're using unsigned integers here, only subtract if the result is > 0
	//	from = mbox.Messages - 3
	//}
	seqset := new(imap.SeqSet)
	seqset.AddRange(1, to)

	messages := make(chan *imap.Message, 10)
	done = make(chan error, 1)
	go func() {
		done <- c.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, messages)
	}()

	log.Printf("messages: %d\n", mbox.Messages)
	//for msg := range messages {
	//log.Printf("* " + msg.Envelope.Subject + "(" + msg.Envelope.Date.String() +")")
	// TODO: handle message -> check date - if not within the last 6 months move to year
	// -> check existence of archive mailbox/year/inbox|outbox
	//}

	//if err := <-done; err != nil {
	//	log.Fatal(err)
	//}

	log.Println("Done!")
	return nil
}
