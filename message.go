//MIMEMail provides convenient formatting (and sending) of MIME formatted emails.
//
//Simply create a new Mail struct with NewMail(), add Recipients (To, Cc, Bcc, etc). Set
//the Subject, add Attachments (by filename), get a Writer for the body
//by calling HTMLBody() or PlainTextBody() and render your template into it.
//Finally call Bytes() to obtain the formatted email or WriteTo() to directly
//write it to a Writer or send it directly (via smtp.SendMail) through the
//Mail.SendMail() method.
package MIMEMail

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"
	"net/smtp"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"
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

// Mail represents a MIME email message and handles encoding,
// MIME headers and so on.
type Mail struct {
	// Address Lists for the Mailheader,
	// Fields that are nil, will be ignored.
	// Use the Add_Recipient or AddAddress for convienice.
	Addresses

	// The subject Line
	Subject string

	parts []*MIMEPart

	msg *bytes.Buffer
}

// Returns a new mail object ready to use.
func NewMail() *Mail {
	return &Mail{
		Addresses: NewAddresses(),
		parts:     make([]*MIMEPart, 0, 1),
		msg:       bytes.NewBuffer(nil),
	}
}

// Sends the mail via smtp.SendMail. If you have the Sender field set, it's first
// entry is used and should match the Address in auth, these values are then passed on to
// smtp.SendMail, returning any errors it throws, else the first From entry
// is used (with the same restrictions). If both are nil, a NoSender error is returned.
func (m *Mail) SendMail(adr string, auth smtp.Auth) error {

	if m.Addresses["Sender"] != nil {
		return smtp.SendMail(adr, auth, m.Addresses["Sender"][0].Address, m.Recipients(), m.msg.Bytes())
	}

	if m.Addresses["From"] != nil {
		return smtp.SendMail(adr, auth, m.Addresses["From"][0].Address, m.Recipients(), m.msg.Bytes())
	}

	return new(NoSender)
}

func (m *Mail) AddFile(filename, attachmentname string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	if attachmentname == "" {
		attachmentname = filepath.Base(filename)
	}

	return m.AddReader(attachmentname, f)
}

func (m *Mail) AddReader(name string, r io.Reader) error {
	// Content-Type: application/octet-stream
	// Content-Transfer-Encoding: base64
	// Content-Disposition: attachment; filename="short_attachment.txt"
	part := NewMIMEPart()
	part.Set(content_type, mime_octetstream)
	part.Set(content_transfer_encoding, mime_base64)
	part.Set(content_disposition, fmt.Sprintf("%s; filename=%s", mime_attachment, name))

	if _, err := io.Copy(base64.NewEncoder(base64.StdEncoding, part.Buffer), r); err != nil {
		return err
	}
	m.parts = append(m.parts, part)
	return nil
}

func (m *Mail) getHeader() textproto.MIMEHeader {
	part := make(textproto.MIMEHeader)

	part = m.ToMimeHeader(part)
	part.Set("Subject", m.Subject)
	part.Set("MIME-Version", "MIME 1.0")
	return part
}

func (m *Mail) addPart(ct, enc string) *MIMEPart {
	n := NewMIMEPart()
	n.Set(content_type, fmt.Sprintf("%s; %s=%s", ct, charset, enc))
	//n.Set(charset, enc)

	m.parts = append(m.parts, n)
	return n
}

// Formats the mail obj for using a HTML body and returns a buffer that you can
// render your Template to. You must call either HTMLBody or PlainTextBody.
// If you call both, only your last call will be respected.
func (m *Mail) HTMLBody() io.Writer {
	return m.addPart(mime_html, mime_utf8)
}

// Formats the mail obj for using a plaintext body and returns a buffer that you
// can render your Template to. You must call either HTMLBody or PlainTextBody.
// If you call both, only your last call will be respected.
func (m *Mail) PlainTextBody() io.Writer {
	return m.addPart(mime_text, mime_utf8)
}

// Returns the fully formatted complete message as a slice of bytes.
// Triggers formatting.
func (m *Mail) Bytes() ([]byte, error) {
	if err := m.write(m.msg); err != nil {
		return nil, err
	}
	return m.msg.Bytes(), nil
}

func (m *Mail) writeHeader(w io.Writer) error {
	header := m.getHeader()
	for field, values := range header {
		if _, err := w.Write([]byte(fmt.Sprintf("%s: %s\r\n", field, strings.Join(values, ", ")))); err != nil {
			return err
		}
	}
	if _, err := w.Write([]byte{13, 10}); err != nil {
		return err
	}
	return nil
}

func (m *Mail) write(w io.Writer) error {
	mpw := multipart.NewWriter(w)

	if err := m.writeHeader(w); err != nil {
		return err
	}

	w.Write([]byte(fmt.Sprintf("%s: %s; boundary=%s\r\n\r\n", content_type, mime_multipart, mpw.Boundary())))

	for _, part := range m.parts {
		pw, err := mpw.CreatePart(part.MIMEHeader)
		if err != nil {
			return err
		}

		if _, err := pw.Write(part.Bytes()); err != nil {
			return err
		}
	}

	if err := mpw.Close(); err != nil {
		return err
	}

	return nil
}

// Writes the fully formatted complete message to the given writer.
// Triggers formatting.
func (m *Mail) WriteTo(w io.Writer) error {
	return m.write(w)
}
