package MIMEMail

import (
	"bytes"
	"net/textproto"
)

type MIMEPart struct {
	textproto.MIMEHeader
	*bytes.Buffer
}

func NewMIMEPart() *MIMEPart {
	return &MIMEPart{
		make(textproto.MIMEHeader),
		bytes.NewBuffer(nil),
	}
}
