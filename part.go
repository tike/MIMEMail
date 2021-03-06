package MIMEMail

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"net/textproto"
	"os"
	"path/filepath"
)

const (
	content_type   = "Content-Type"
	charset        = "charset"
	mime_multipart = "multipart/mixed"
	mime_html      = "text/html"
	mime_text      = "text/plain"
	mime_utf8      = "utf-8"

	mime_octetstream          = "application/octet-stream"
	content_transfer_encoding = "Content-Transfer-Encoding"
	mime_base64               = "base64"

	content_disposition = "Content-Disposition"
	mime_attachment     = "attachment"
)

// MIMEPart wraps the MIMEPart functionality, in all likelyhood you'll never use
// it directly.
type MIMEPart struct {
	textproto.MIMEHeader
	*bytes.Buffer
}

// NewMIMEPart creates a new blank MIMEPart.
func NewMIMEPart() *MIMEPart {
	return &MIMEPart{
		make(textproto.MIMEHeader),
		bytes.NewBuffer(nil),
	}
}

// NewPart creates a new MIMEPart with the given Content-Type and encoding.
func NewPart(contenttype, encoding string) *MIMEPart {
	p := NewMIMEPart()
	p.Set(content_type, fmt.Sprintf("%s; charset=%s", contenttype, encoding))
	return p
}

// NewHTML creates a new MIMEPart with "Content-Type: text/html; encoding: utf-8"
func NewHTML() *MIMEPart {
	return NewPart(mime_html, mime_utf8)
}

// NewPlainText creates a new MIMEPart with "Content-Type: text/plain; encoding: utf-8"
func NewPlainText() *MIMEPart {
	return NewPart(mime_text, mime_utf8)
}


// NewPGPVersion creates a new PGP/MIME Version header.
func NewPGPVersion() *MIMEPart {
	p := NewMIMEPart()
	p.Set(content_type, "application/pgp-encrypted")
	p.Set(content_disposition, "PGP/MIME version identification")
	p.WriteString("Version: 1\r\n")
	return p
}

// NewPGPBody creates a PGP/MIME message body part.
func NewPGPBody() *MIMEPart {
	p := NewMIMEPart()
	p.Set(content_type, "application/pgp-encrypted")
	p.Set(content_disposition, `inline; filename="encrypted.asc"`)
	return p
}

// NewAttachment creates a new MIMEPart with all the necessary headers set.
func NewAttachment(name string, r io.Reader) (*MIMEPart, error) {
	// Content-Type: application/octet-stream
	// Content-Transfer-Encoding: base64
	// Content-Disposition: attachment; filename="short_attachment.txt"
	p := NewMIMEPart()
	p.Set(content_type, mime_octetstream)
	p.Set(content_transfer_encoding, mime_base64)
	p.Set(content_disposition, fmt.Sprintf("%s; filename=%s", mime_attachment, name))

	if _, err := io.Copy(base64.NewEncoder(base64.StdEncoding, p.Buffer), r); err != nil {
		return nil, err
	}
	return p, nil
}

// NewFile creates a new File attachment MIMEPart with all the necessary headers set.
// If you pass a string as the optional attachment argument, it will be used as the
// filename for sending the attachment, if no such argument is passed, filepath.Base(file)
// will be used.
func NewFile(file string, attachment ...string) (*MIMEPart, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	attachmentname := filepath.Base(file)
	if len(attachment) != 0 {
		attachmentname = attachment[0]
	}

	return NewAttachment(attachmentname, f)
}
