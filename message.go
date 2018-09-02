package MIMEMail

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/smtp"
	"net/textproto"

	"github.com/tike/MIMEMail/templated"
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

	// for testing purposes only
	boundary string
}

// NewMail returns a new mail object ready to use.
func NewMail() *Mail {
	return &Mail{
		Addresses: NewAddresses(),
		parts:     make([]*MIMEPart, 0, 1),
	}
}

// NewTemplated renders the template specified by config and name using data as the
// rendering context. The results will be put into the Subject and HTMLBody of the
//  returned Mail struct.
func NewTemplated(cnf *templated.Config, name string, data interface{}) (*Mail, error) {
	subj, body, err := templated.Render(cnf, name, data)
	if err != nil {
		return nil, err
	}

	m := NewMail()
	m.Subject = subj

	bodyPart := NewHTML()
	bodyPart.Buffer = bytes.NewBuffer(body)
	m.parts = append(m.parts, bodyPart)

	return m, nil
}

// SendMail sends the mail via smtp.SendMail (which uses StartTLS if available).
// If you have the Sender field set, it's first entry is used and
// should match the Address in auth, else the first "From" entry
// is used (with the same restrictions). If both are nil,
// a NoSender error is returned.
// These values are then passed on to smtp.SendMail, returning any errors it throws.
func (m *Mail) SendMail(adr string, auth smtp.Auth) error {
	msg, err := m.Bytes()
	if err != nil {
		return err
	}

	from, err := m.EffectiveSender()
	if err != nil {
		return err
	}

	return smtp.SendMail(adr, auth, from, m.Recipients(), msg)
}

// AddFile adds the file given by filename as an attachment to the mail.
// If you provide the optional attachmentname argument, the file will be
// attached with this name.
func (m *Mail) AddFile(filename string, attachmentname ...string) error {
	p, err := NewFile(filename, attachmentname...)
	if err != nil {
		return err
	}

	m.parts = append(m.parts, p)
	return nil
}

// AddReader adds the given reader as an attachment, using name as the filename.
func (m *Mail) AddReader(name string, r io.Reader) error {
	p, err := NewAttachment(name, r)
	if err != nil {
		return err
	}

	m.parts = append(m.parts, p)
	return nil
}

func (m *Mail) getHeader() textproto.MIMEHeader {
	part := make(textproto.MIMEHeader)

	part = m.ToMimeHeader(part)
	part.Set("Subject", m.Subject)
	part.Set("MIME-Version", "1.0")
	return part
}

// HTMLBody adds a HTML body part and returns a buffer that you can render your Template to.
func (m *Mail) HTMLBody() io.Writer {
	p := NewHTML()
	m.parts = append(m.parts, p)
	return p
}

// PlainTextBody adds a Plaintext body part and returns a buffer that you can render your Template to.
func (m *Mail) PlainTextBody() io.Writer {
	p := NewPlainText()
	m.parts = append(m.parts, p)
	return p
}

// Bytes returns the fully formatted complete message as a slice of bytes.
// This is for plain MIME mails, if you want a PGP/MIME encrypted mail, use the Encrypt method instead.
func (m *Mail) Bytes() ([]byte, error) {
	msg := bytes.NewBuffer(nil)
	if err := m.write(msg); err != nil {
		return nil, err
	}
	return msg.Bytes(), nil
}

var headerOrder = []string{"Sender", "From", "To", "Cc", "Bcc", "ReplyTo", "FollowupTo", "Subject", "MIME-Version"}

func (m *Mail) writeHeader(w io.Writer) error {
	header := m.getHeader()
	for _, field := range headerOrder {
		value := header.Get(field)
		if value == "" {
			continue
		}

		if _, err := w.Write([]byte(fmt.Sprintf("%s: %s\r\n", field, value))); err != nil {
			return err
		}
	}

	return nil
}

func (m *Mail) writeBody(w io.Writer) error {
	mpw := multipart.NewWriter(w)
	m.boundary = mpw.Boundary()

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

	return mpw.Close()
}

func (m *Mail) write(w io.Writer) error {
	if err := m.writeHeader(w); err != nil {
		return err
	}

	if err := m.writeBody(w); err != nil {
		return err
	}

	return nil
}

// WriteTo writes the fully formatted complete message to the given writer.
// This is for plain MIME mails, if you want a PGP/MIME encrypted mail, use the WriteEncrypted method instead.
func (m *Mail) WriteTo(w io.Writer) error {
	return m.write(w)
}

// Encrypt encrypts the mail with PGP/MIME using CreateEntity to obtain recpient and CreateSigningEntity to obtain the signing entity.
// If signer is nil, the mail will simply not be signed. The Key Fields of both to and signer must be non-nil.
func (m *Mail) Encrypt(to *Account, signer *Account) ([]byte, error) {
	var b bytes.Buffer
	if err := m.WriteEncrypted(&b, to, signer); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

// WriteEncrypted encrypts the mail with PGP/MIME using CreateEntity to obtain the recpient and CreateSigningEntity to obtain the signing entity.
// If signer is nil, the mail will simply not be signed. The Key Fields of both to and signer must be non-nil.
func (m *Mail) WriteEncrypted(w io.Writer, to *Account, signer *Account) error {
	if err := m.writeHeader(w); err != nil {
		return err
	}

	mpw := multipart.NewWriter(w)
	pgpMIMEheader := fmt.Sprintf("%s: %s; protocol=%q; boundary=%q\r\n\r\n",
		content_type, "multipart/encrypted", "application/pgp-encrypted", mpw.Boundary())
	if _, err := w.Write([]byte(pgpMIMEheader)); err != nil {
		return err
	}

	pgpVersion := NewPGPVersion()
	pgpVersionBody, err := mpw.CreatePart(pgpVersion.MIMEHeader)
	if err != nil {
		return err
	}
	if _, err := pgpVersionBody.Write(pgpVersion.Bytes()); err != nil {
		return err
	}

	pgpBody := NewPGPBody()
	pgpBodyPart, err := mpw.CreatePart(pgpBody.MIMEHeader)
	if err != nil {
		return err
	}

	plainTextWriter, err := Encrypt(pgpBodyPart, to, signer)
	if err != nil {
		return err
	}

	if err := m.writeBody(plainTextWriter); err != nil {
		return err
	}

	if err := plainTextWriter.Close(); err != nil {
		return err
	}
	if _, err := pgpBodyPart.Write([]byte("\r\n")); err != nil {
		return err
	}

	return mpw.Close()
}
