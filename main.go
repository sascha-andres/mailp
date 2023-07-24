package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/sascha-andres/mailp/internal"

	"github.com/sascha-andres/mailp/internal/config"

	"github.com/sascha-andres/reuse/flag"
)

var (
	configPath, folder   string
	mailID, outputFormat string
	debug                bool
)

const (
	// Prefix is logging prefix and env prefix
	Prefix = "MAILP"
)

func init() {
	log.SetPrefix(fmt.Sprintf("[%s] ", Prefix))
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.LUTC)

	flag.SetEnvPrefix(Prefix)
	flag.StringVar(&outputFormat, "output", "json", "output format")
	flag.StringVar(&configPath, "config", "", "path to config file")
	flag.StringVarWithoutEnv(&folder, "folder", "INBOX", "list mails from folder")
	flag.StringVarWithoutEnv(&mailID, "mail", "", "show mail with id")
	flag.BoolVar(&debug, "debug", false, "enable debug output")
}

func main() {
	flag.Parse()
	if debug {
		log.Println("Starting mailp...")
	}

	if configPath == "" {
		log.Fatal("no config file specified")
	}

	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	if debug {
		log.Print("reading config from ", configPath)
	}
	cfg, err := config.FromFile(configPath)
	if err != nil {
		return err
	}
	connector, err := internal.NewConnector("imap")
	if err != nil {
		return err
	}
	if err := connector.Initialize(cfg); err != nil {
		return err
	}

	verbs := flag.GetVerbs()
	if len(verbs) == 0 {
		return errors.New("no verb specified. expected: list, show, folder")
	}

	if verbs[0] == "list" {
		return listMails(connector)
	}

	if verbs[0] == "show" {
		return showMail(connector)
	}

	if verbs[0] == "folder" {
		return listFolder(connector)
	}

	return errors.New(fmt.Sprintf("unknown verb: %q", verbs[0]))
}

func listFolder(connector internal.Connector) error {
	folders, err := connector.ListFolder()
	if err != nil {
		return err
	}
	if outputFormat == "json" {
		data, err := json.Marshal(folders)
		if err != nil {
			return err
		}
		fmt.Printf("%s\n", data)
	}
	if outputFormat == "text" {
		for _, folder := range folders {
			fmt.Println(folder)
		}
	}
	return nil
}

func listMails(connector internal.Connector) error {
	if folder == "" {
		return errors.New("no folder specified")
	}
	mails, err := connector.ListMails(folder)
	if err != nil {
		return err
	}
	if outputFormat == "json" {
		data, err := json.Marshal(mails)
		if err != nil {
			return err
		}
		fmt.Printf("%s\n", data)
	}
	if outputFormat == "text" {
		tw := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
		for _, mail := range mails {
			_, _ = tw.Write([]byte(fmt.Sprintf("%s\t%s\t%s\t%s\n", mail.ID, mail.Received, mail.From, mail.Subject)))
		}
		_ = tw.Flush()
	}
	return nil
}

func showMail(connector internal.Connector) error {
	if mailID == "" {
		return errors.New("no mail id specified")
	}
	if folder == "" {
		return errors.New("no folder specified")
	}
	if debug {
		log.Print("showing mail ", mailID)
	}
	m, err := connector.GetMail(folder, mailID)
	if err != nil {
		return err
	}
	if outputFormat == "json" {
		data, err := json.Marshal(m)
		if err != nil {
			return err
		}
		fmt.Printf("%s\n", data)
	}
	if outputFormat == "text" {
		to := strings.Join(m.To, ", ")
		cc := strings.Join(m.Cc, ", ")
		tw := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
		_, _ = tw.Write([]byte(fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s\n%s", m.ID, to, cc, m.Received, m.From, m.Subject, m.Body)))
	}
	return nil
}
