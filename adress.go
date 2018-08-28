package MIMEMail

import (
	"net/mail"
	"net/textproto"
)

// AddressHeader is a dedicated type for mail AddressHeader fields.
type AddressHeader string

// Valid values for mail address headers
const (
	AddrSender     AddressHeader = "Sender"
	AddrFrom       AddressHeader = "From"
	AddrTo         AddressHeader = "To"
	AddrCc         AddressHeader = "Cc"
	AddrBcc        AddressHeader = "Bcc"
	AddrReplyTo    AddressHeader = "ReplyTo"
	AddrFollowupTo AddressHeader = "FollowupTo"
)

// Addresses handles setting and encoding the mail address headers
type Addresses map[AddressHeader][]mail.Address

// NewAddresses creates a new mail address header
func NewAddresses() Addresses {
	return make(Addresses, 3)
}

// Recipients returns just the mailaddresses of all the recipients
// (To, Cc, Bcc), ready to be passed to smtp.SendMail et al. for your convenience.
func (a Addresses) Recipients() []string {
	to := make([]string, 0, 10)
	for _, field := range []AddressHeader{AddrTo, AddrCc, AddrBcc} {
		if a[field] != nil {
			for _, address := range a[field] {
				to = append(to, address.Address)
			}
		}
	}
	return to
}

// Sender adds the given name, address pair to Sender.
func (a *Addresses) Sender(name, address string) error {
	return a.AddPerson(AddrSender, name, address)
}

// SenderAddr adds the given address to Sender.
func (a *Addresses) SenderAddr(address mail.Address) error {
	return a.AddAddress(AddrSender, address)
}

// From adds the given name, address pair to From.
func (a *Addresses) From(name, address string) error {
	return a.AddPerson(AddrFrom, name, address)
}

// FromAddr adds the given address to From.
func (a *Addresses) FromAddr(address mail.Address) error {
	return a.AddAddress(AddrFrom, address)
}

// To adds the given name, address pair to To.
func (a *Addresses) To(name, address string) error {
	return a.AddPerson(AddrTo, name, address)
}

// ToAddr adds the given address to To.
func (a *Addresses) ToAddr(address mail.Address) error {
	return a.AddAddress(AddrTo, address)
}

// Cc adds the given name, address pair to Cc.
func (a *Addresses) Cc(name, address string) error {
	return a.AddPerson(AddrCc, name, address)
}

// CcAddr adds the given address to Cc.
func (a *Addresses) CcAddr(address mail.Address) error {
	return a.AddAddress(AddrCc, address)
}

// Bcc adds the given name, address pair to Bcc.
func (a *Addresses) Bcc(name, address string) error {
	return a.AddPerson(AddrBcc, name, address)
}

// BccAddr adds the given address to Bcc.
func (a *Addresses) BccAddr(address mail.Address) error {
	return a.AddAddress(AddrBcc, address)
}

// ReplyTo adds the given name, address pair to ReplyTo.
func (a *Addresses) ReplyTo(name, address string) error {
	return a.AddPerson(AddrReplyTo, name, address)
}

// ReplyToAddr adds the given address to ReplyTo.
func (a *Addresses) ReplyToAddr(address mail.Address) error {
	return a.AddAddress(AddrReplyTo, address)
}

// FollowupTo adds the given name, address pair to FollowupTo.
func (a *Addresses) FollowupTo(name, address string) error {
	return a.AddPerson(AddrFollowupTo, name, address)
}

// FollowupToAddr adds the given address to FollowupTo.
func (a *Addresses) FollowupToAddr(address mail.Address) error {
	return a.AddAddress(AddrFollowupTo, address)
}

// AddPerson adds the given details to the given mail header field.
// Field should be a valid address field. Use the predefined Addr... constants.
// Adding the will fail if the field is not one of the predefined constants.
func (a *Addresses) AddPerson(field AddressHeader, name, address string) error {
	return a.AddAddress(field, mail.Address{Name: name, Address: address})
}

// AddAddress a recipient to your mail Header.
// Field should be a valid address field. Use the predefined Addr... constants.
// Adding the will fail if the field is not one of the predefined constants.
func (a *Addresses) AddAddress(field AddressHeader, address mail.Address) error {
	if !valid(field) {
		return InvalidField(field)
	}

	(*a)[field] = append((*a)[field], address)
	return nil
}

func valid(field AddressHeader) bool {
	switch field {
	case AddrSender, AddrFrom, AddrTo, AddrCc, AddrBcc, AddrReplyTo, AddrFollowupTo:
		return true
	default:
		return false
	}
}

// ToMimeHeader packs the contents up in the given MIMEHeader or creates a new
// one if nil is passed.
func (a Addresses) ToMimeHeader(part textproto.MIMEHeader) textproto.MIMEHeader {
	if part == nil {
		part = make(textproto.MIMEHeader)
	}

	for field, addresses := range a {
		if addresses != nil {
			for _, address := range addresses {
				part.Add(string(field), address.String())
			}
		}
	}
	return part
}
