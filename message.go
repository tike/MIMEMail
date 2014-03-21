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
	"io"
	"mime/multipart"
	"net/smtp"
	"net/textproto"
	"os"
)

const content_type = "Content-Type"
const charset = "charset"
const mime_html = "text/html"
const mime_text = "text/plain"
const mime_utf8 = "utf-8"

const MULTI = "Content-Type: multipart/mixed; boundary="

const HTML = "Content-Type: text/html; charset=utf-8"
const TEXT = "Content-Type: text/plain; charset=utf-8"

const FILE = "Content-Type: application/octet-stream"
const FILE_TFE = "Content-Transfer-Encoding: base64"
const FILE_DISP = "Content-Disposition: attachment; filename="

// Mail represents a MIME email message and handles encoding,
// MIME headers and so on.
type Mail struct {
	// Address Lists for the Mailheader,
	// Fields that are nil, will be ignored.
	// Use the Add_Recipient or AddAddress for convienice.
	Addresses

	// The subject Line
	Subject string

	parts       []*MIMEPart
	attachments []io.Reader

	msg *bytes.Buffer
}

// Returns a new mail object ready to use.
func NewMail() *Mail {
	return &Mail{
		Addresses:   NewAddresses(),
		parts:       make([]*MIMEPart, 0, 1),
		attachments: make([]io.Reader, 0, 1),
		msg:         bytes.NewBuffer(nil),
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

func (m *Mail) AddFile(filename string) error {
	r, err := os.Open(filename)
	if err != nil {
		return err
	}
	m.attachments = append(m.attachments, r)
	return nil
}

func (m *Mail) AddReader(r io.Reader) {
	m.attachments = append(m.attachments, r)
}

func (m *Mail) getHeader() textproto.MIMEHeader {
	part := make(textproto.MIMEHeader)
	part.Set("MIME-Version", "MIME 1.0")
	part.Set("Subject", m.Subject)
	part = m.ToMimeHeader(part)
	return part
}

func (m *Mail) addPart(ct, enc string) *MIMEPart {
	n := NewMIMEPart()
	n.Set(content_type, ct)
	n.Set(charset, enc)

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

func (m *Mail) write(w io.Writer) error {
	mpw := multipart.NewWriter(w)
	if _, err := mpw.CreatePart(m.getHeader()); err != nil {
		return err
	}

	for _, part := range m.parts {
		pw, err := mpw.CreatePart(part.MIMEHeader)
		if err != nil {
			return err
		}

		if _, err := pw.Write(part.Bytes()); err != nil {
			return err
		}
	}

	for _, att := range m.attachments {
		pw, err := mpw.CreateFormFile("foo", "bar")
		if err != nil {
			return err
		}
		if _, err = io.Copy(pw, att); err != nil {
			return err
		}
	}

	return nil
}

// Writes the fully formatted complete message to the given writer.
// Triggers formatting.
func (m *Mail) WriteTo(w io.Writer) error {
	return m.write(w)
}
