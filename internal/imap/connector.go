package imap

import (
	"errors"
	"fmt"
	"io"
	"log"

	"github.com/sascha-andres/mailp/internal/data"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/mail"
	"github.com/sascha-andres/mailp/internal/config"

	_ "github.com/emersion/go-message/charset"
)

type Connector struct {
	client *client.Client
}

func (c *Connector) Initialize(cfg *config.Config) error {
	i, err := client.DialTLS(fmt.Sprintf("%s:%d", cfg.Imap.Host, cfg.Imap.Port), nil)
	if err != nil {
		return err
	}
	err = i.Login(cfg.Imap.Username, cfg.Imap.Password)
	if err != nil {
		return err
	}
	c.client = i
	return nil
}

func (c *Connector) ListMails(folder string) ([]*data.Mail, error) {
	result := make([]*data.Mail, 0)

	mbox, err := c.client.Select(folder, true)
	if err != nil {
		return nil, err
	}
	if mbox.Messages == 0 {
		return result, nil
	}

	seqSet := new(imap.SeqSet)
	seqSet.AddRange(1, mbox.Messages)

	// Get the whole message body
	var section imap.BodySectionName
	items := []imap.FetchItem{imap.FetchUid, section.FetchItem()} // section.FetchItem(),

	messages := make(chan *imap.Message, 1)
	done := make(chan error, 1)
	go func() {
		done <- c.client.Fetch(seqSet, items, messages)
	}()

	for {
		msg := <-messages

		if msg == nil {
			break
		}

		m := data.Mail{}

		r := msg.GetBody(&section)
		if r == nil {
			log.Fatal("Server didn't returned message body")
		}

		mr, err := mail.CreateReader(r)
		if err != nil {
			log.Fatal(err)
		}

		header := mr.Header
		if date, err := header.Date(); err == nil {
			m.Received = date.Format("2006-01-02 15:04:05")
		}
		if from, err := header.AddressList("From"); err == nil {
			m.From = from[0].String()
		}
		if to, err := header.AddressList("To"); err == nil {
			m.To = make([]string, 0)
			for _, t := range to {
				m.To = append(m.To, t.String())
			}
		}
		if cc, err := header.AddressList("Cc"); err == nil {
			m.Cc = make([]string, 0)
			for _, t := range cc {
				m.Cc = append(m.To, t.String())
			}
		}
		if subject, err := header.Subject(); err == nil {
			m.Subject = subject
		}
		m.ID = fmt.Sprintf("%d", msg.Uid)

		result = append(result, &m)
	}

	return result, <-done
}

func (c *Connector) GetMail(folder, id string) (*data.MailData, error) {
	_, err := c.client.Select(folder, true)
	if err != nil {
		return nil, err
	}

	seqSet := new(imap.SeqSet)
	err = seqSet.Add(id)
	if err != nil {
		return nil, err
	}

	var section imap.BodySectionName
	items := []imap.FetchItem{section.FetchItem(), imap.FetchUid}

	messages := make(chan *imap.Message, 1)
	done := make(chan error, 1)
	go func() {
		done <- c.client.UidFetch(seqSet, items, messages)
	}()

	var m data.MailData

	msg := <-messages

	r := msg.GetBody(&section)
	if r == nil {
		return nil, errors.New("server didn't returned message body")
	}

	mr, err := mail.CreateReader(r)
	if err != nil {
		return nil, err
	}

	header := mr.Header
	if date, err := header.Date(); err == nil {
		m.Received = date.Format("2006-01-02 15:04:05")
	}
	if from, err := header.AddressList("From"); err == nil {
		m.From = from[0].String()
	}
	if to, err := header.AddressList("To"); err == nil {
		m.To = make([]string, 0)
		for _, t := range to {
			m.To = append(m.To, t.String())
		}
	}
	if cc, err := header.AddressList("Cc"); err == nil {
		m.Cc = make([]string, 0)
		for _, t := range cc {
			m.Cc = append(m.To, t.String())
		}
	}
	if subject, err := header.Subject(); err == nil {
		m.Subject = subject
	}
	m.ID = fmt.Sprintf("%d", msg.Uid)
	m.Attachments = make([]data.Attachment, 0)

	// Process each message's part
	for {
		p, err := mr.NextPart()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		switch h := p.Header.(type) {
		case *mail.InlineHeader:
			t, _, err := h.ContentType()
			if err != nil {
				log.Printf("error reading content-type: %s", err)
				continue
			}
			if t == "text/plain" {
				b, _ := io.ReadAll(p.Body)
				m.Body = string(b)
			}
		case *mail.AttachmentHeader:
			// This is an attachment
			filename, _ := h.Filename()
			m.Attachments = append(m.Attachments, data.Attachment{Name: filename})
		}
	}

	return &m, nil
}

func (c *Connector) ListFolder() ([]string, error) {
	result := make([]string, 1)
	result[0] = "INBOX"

	mailboxes := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)
	go func() {
		done <- c.client.List("", "*", mailboxes)
	}()

	for m := range mailboxes {
		result = append(result, m.Name)
	}

	return result, <-done
}
