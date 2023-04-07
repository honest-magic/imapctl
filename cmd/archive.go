package cmd

import (
	"bitbucket.org/mis79/imapctl/utl"
	"fmt"
	imap "github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/spf13/cobra"
	"log"
	"strconv"
	"time"
)

var archiveCmd = &cobra.Command{
	Use:   "archive",
	Short: "Archives mail messages into an archive folder",
	Long:  "Archives mail messages into an archive folder",
	RunE:  cmdArchive,
}

type ArchiveState struct {
	Name  string
	Count int
	Seq   *imap.SeqSet
}

var target string
var inbox string
var outbox []string

var targetInbox string
var targetOutbox string

var location string
var dryrun bool

func init() {

	archiveCmd.Flags().StringVarP(&target, "target", "t", "INBOX.Archive", "The target mailbox for the archive")
	archiveCmd.Flags().StringVarP(&inbox, "inbox", "i", "INBOX", "The mailbox name of received messages")
	// Instead we could scan for the '\Sent' flag
	archiveCmd.Flags().StringSliceVarP(&outbox, "outbox", "o", []string{"INBOX.Sent"}, "The mailbox name of sent messages")
	archiveCmd.Flags().StringVarP(&targetInbox, "target-inbox", "r", "Inbox", "The archive mailbox name of received messages")
	archiveCmd.Flags().StringVarP(&targetOutbox, "target-outbox", "l", "Sent", "The archive mailbox name of sent messages")
	archiveCmd.Flags().StringVarP(&location, "location", "z", "Europe/Zurich", "The timezone to use")
	archiveCmd.Flags().BoolVarP(&dryrun, "dryrun", "d", false, "Flag to control whether the command is just a dry run to check the outcome with verbose")

	rootCmd.AddCommand(archiveCmd)
}

func cmdArchive(cmd *cobra.Command, args []string) error {

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

	// get the current year
	loc, _ := time.LoadLocation(location)
	now := time.Now().In(loc)
	//year := now.Year()

	// Select INBOX
	mbox, err := c.Select(inbox, true)
	if err != nil {
		return err
	}

	to := mbox.Messages
	sequentialSet := new(imap.SeqSet)
	sequentialSet.AddRange(1, to)

	messages := make(chan *imap.Message, 10)
	done := make(chan error, 1)
	go func() {
		done <- c.Fetch(sequentialSet, []imap.FetchItem{imap.FetchEnvelope, imap.FetchUid, imap.FetchInternalDate}, messages)
	}()

	moveMap := make(map[string]*ArchiveState)
	limit := now.AddDate(0, -6, 0)
	for msg := range messages {

		msgDate := msg.Envelope.Date.In(loc)
		msgIntDate := msg.InternalDate.In(loc)

		if msgIntDate.Before(limit) {
			boxName := target + "." + strconv.Itoa(msgIntDate.Year()) + "." + targetInbox

			state := moveMap[boxName]
			if state == nil {
				state = &ArchiveState{boxName, 0, new(imap.SeqSet)}
				moveMap[boxName] = state
				if verbose {
					log.Println("Created entry for: " + boxName)
				}
			}

			state.Seq.AddNum(msg.Uid)
			state.Count++
			if verbose {
				uid := msg.Uid
				log.Printf("Added %s '%s' (date: %s, internal: %s,  uid: %s) to %s\n", msg.SeqNum, msg.Envelope.Subject, msgDate, msgIntDate, uid, boxName)
			}

		}
	}

	if err := <-done; err != nil {
		return err
	}

	for k, v := range moveMap {

		if verbose {
			log.Printf("%s (%d) -> %s", v.Seq, v.Count, k)
		}

		if !dryrun {
			if err := utl.CreateDirectory(c, k); err != nil {
				return err
			}
			if err := c.UidMove(v.Seq, k); err != nil {
				return err
			}
		}
	}

	if verbose {
		log.Println("Dry run!")
	}

	return nil
}
