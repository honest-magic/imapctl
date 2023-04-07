package utl

import (
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"log"
)

func CreateDirectory(c *client.Client, name string) error {

	boxes := make(chan *imap.MailboxInfo, 1)
	done := make(chan error, 1)
	go func() {
		done <- c.List("", name, boxes)
	}()

	found := false

	for _ = range boxes {
		found = true
	}

	if err := <-done; err != nil {
		return err
	}

	if !found {
		log.Println("Create directory: " + name)
		if err := c.Create(name); err != nil {
			return err
		}
	}

	return nil
}
