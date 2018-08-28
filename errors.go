package MIMEMail

type NoSender int

func (e NoSender) Error() string {
	return "You have neither From nor Sender set on your Mail!"
}

type InvalidField AddressHeader

func (e InvalidField) Error() string {
	return string(e) + " is not a valid field (use: From, Sender, To, Cc, Bcc, ReplyTo or FollowupTo)"
}
