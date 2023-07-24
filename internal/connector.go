package internal

import (
	"errors"

	"github.com/sascha-andres/mailp/internal/data"

	"github.com/sascha-andres/mailp/internal/config"
	"github.com/sascha-andres/mailp/internal/imap"
)

// Connector is the interface for connectors
type (
	Connector interface {
		//Initialize initializes the connector with the configuration
		Initialize(cfg *config.Config) error
		// ListMails lists all folders available
		ListMails(folder string) ([]*data.Mail, error)
		// GetMail returns the mail with the given id
		GetMail(folder, id string) (*data.MailData, error)
		// ListFolder lists all folders available
		ListFolder() ([]string, error)
	}
)

// NewConnector returns a new connector
func NewConnector(t string) (Connector, error) {
	switch t {
	case "imap":
		return &imap.Connector{}, nil
	}
	return nil, errors.New("unknown connector type")
}
