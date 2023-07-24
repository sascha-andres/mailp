# mail printer

list mails in IMAP folder and print data on console.

## Confifguaration

Configuration for mail server can be done in `~/.config/mailp/config.json`:

    {
        "host": "imap.example.com",
        "port": 993,
        "username": "<>",
        "password": "<>"
    }

The environment variable `MAILP_CONFIG` can be used to set a path to the configuration file. Alternatively the path can be set with command line flags.

## Flags

| Flag    | Environment variable | Description                                                                   | Default |
| ------- | -------------------- | ----------------------------------------------------------------------------- | ------- |
| -output | MAILP_OUTPUT         | json or text, json will print a valid json document, text tab separated lines | json    |
| -folder |                      | IMAP folder to list, default is INBOX                                         | INBOX   |
| -config | MAILP_CONFIG         | path to config file, default is ~/.config/mailp/config.json                   |         |
| -debug  | MAILP_DEBUG          | print debug information                                                       | false   |
| -mail   |                      | print mail content for message with uid                                       |         |
## Usage

### list folders

    mailp folder

Will print a list of all folders in imap account.

### list mails

    mailp list -folder INBOX

If folder is empty an error will be returned.

### print mail content

    mailp -mail 1234 -folder INBOX

Print mail.