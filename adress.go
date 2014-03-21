package MIMEMail

import (
	"net/mail"
	"net/textproto"
)

type Addresses map[string][]mail.Address

func NewAddresses() Addresses {
	return Addresses{
		"Sender":     nil,
		"From":       nil,
		"To":         nil,
		"Cc":         nil,
		"Bcc":        nil,
		"ReplyTo":    nil,
		"FollowupTo": nil,
	}
}

// Returns just the mailaddresses of all the recipients (To, Cc, Bcc), ready to be passed to
// smtp.SendMail et al. for your convenience.
func (a Addresses) Recipients() []string {
	to := make([]string, 0, 10)
	for _, field := range []string{"To", "Cc", "Bcc"} {
		if a[field] != nil {
			for _, address := range a[field] {
				to = append(to, address.Address)
			}
		}
	}
	return to
}

// Adds a recipient to your mail Header. Field should be any of the header address-list fields, i.e.
// "Sender", "From", "To", "Cc", "Bcc", "ReplyTo" or "FollowupTo", otherwise
// adding will fail (return false). Name should be the Name to display and address the email address.
func (a *Addresses) AddPerson(field, name, address string) error {
	return a.AddAddress(field, mail.Address{name, address})
}

// Adds a recipient to your mail Header. Field should be any of the header address-list fields, i.e.
// "Sender", "From", "To", "Cc", "Bcc", "ReplyTo" or "FollowupTo", otherwise
// adding will fail (return false).
// see net/mail for details on the mail.Address struct.
func (a *Addresses) AddAddress(field string, address mail.Address) error {
	_, validField := (*a)[field]
	if !validField {
		return InvalidField(field)
	}
	(*a)[field] = append((*a)[field], address)
	return nil
}

func (a Addresses) ToMimeHeader(part textproto.MIMEHeader) textproto.MIMEHeader {
	if part == nil {
		part = make(textproto.MIMEHeader)
	}

	for field, addresses := range a {
		if addresses != nil {
			for _, address := range addresses {
				part.Add(field, address.String())
			}
		}
	}
	return part
}
