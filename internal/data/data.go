package data

type (
	// Mail is the mail structure
	Mail struct {
		// ID is the mail id
		ID string
		// Subject is the mail subject
		Subject string
		// From is the mail from
		From string
		// To is the mail to
		To []string
		// Cc is the mail cc
		Cc []string
		// Received is the mail received date
		Received string
	}

	// MailData is the mail data structure
	// it embeds the mail structure and additional data
	MailData struct {
		Mail
		// Body is the mail body
		Body string
		// Header is the mail header
		Header []string
		//	Attachments is the list of attachments
		Attachments []Attachment
	}

	// Attachment is the attachment structure
	Attachment struct {
		// Name is the attachment name
		Name string
	}
)
