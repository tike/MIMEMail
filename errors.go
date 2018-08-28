package MIMEMail

// NoSender is returned when trying to send a mail with no From or Sender
// address set.
type NoSender int

func (e NoSender) Error() string {
	return "You have neither From nor Sender set on your Mail!"
}

// InvalidField is returned by the various Addresses' From(Addr), To(Addr), ...
// methods if an invalid field value is provided.
type InvalidField AddressHeader

func (e InvalidField) Error() string {
	return string(e) + " is not a valid field (use: From, Sender, To, Cc, Bcc, ReplyTo or FollowupTo)"
}
