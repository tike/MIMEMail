//MIMEMail provides convenient formatting (and sending) of MIME formatted emails.
//
//Simply instanciate the Mail struct, add Recipients (To, Cc, Bcc). Set Reply-To
//etc, set the Subject, add Attachments (by filename), get a Writer for the body
//by calling HTMLBody() or PlainTextBody() and render your template into it.
//Finally call Bytes() to obtain the formatted email or WriteTo() to directly
//write it to a Writer or send it directly (via smtp.SendMail) through the
//Mail.SendMail() method.
package MIMEMail

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"net/mail"
	"net/smtp"
	"os"
	"path/filepath"
	"strings"
)

//Mail represents a MIME email message and handles encoding.
type Mail struct {
	// Address Lists for the Mailheader,
	// Fields that are nil, will be ignored.
	// Use the Add_Recipient or Add_Address for convienice.
	Recv map[string][]mail.Address

	// The subject Line
	Subject string

	// Filenames of Attachments to send along
	// ignored if nil (default).
	Attachments []string

	boundary   []byte
	mimeHeader string
	bodyHeader string

	body *bytes.Buffer
	out  *bytes.Buffer
}

// Returns a new mail object ready to use.
func NewMail() *Mail {
	return &Mail{
		Recv: map[string][]mail.Address{
			"Sender":     nil,
			"From":       nil,
			"To":         nil,
			"Cc":         nil,
			"Bcc":        nil,
			"ReplyTo":    nil,
			"FollowupTo": nil,
		},
		body: bytes.NewBuffer(nil),
		out:  bytes.NewBuffer(nil),
	}
}

// Adds a recipient to your mail Header. Field should be any of the header address-list fields, i.e.
// "Sender", "From", "To", "Cc", "Bcc", "ReplyTo" or "FollowupTo", otherwise
// adding will fail (return false). Name should be the Name to display and address the
// email address.
func (m *Mail) Add_Recipient(field, name, address string) bool {
	return m.Add_Address(field, mail.Address{name, address})
}

// Adds a recipient to your mail Header. Field should be any of the header address-list fields, i.e.
// "Sender", "From", "To", "Cc", "Bcc", "ReplyTo" or "FollowupTo", otherwise
// adding will fail (return false).
// see net/mail for details on the mail.Address struct.
func (m *Mail) Add_Address(field string, address mail.Address) (added bool) {
	_, validField := m.Recv[field]
	if !validField {
		return false
	}
	m.Recv[field] = append(m.Recv[field], address)
	return true
}

// Returns just the mailaddresses of all the recipients, ready to be passed to
// smtp.SendMail et al. for your convenience.
func (m *Mail) Recipients() (to []string) {
	to = make([]string, 0, 10)
	for _, field := range []string{"To", "Cc", "Bcc"} {
		if m.Recv[field] != nil {
			for _, address := range m.Recv[field] {
				to = append(to, address.Address)
			}
		}
	}
	return
}

// Sends the mail via smtp.SendMail. The first entry of the FROM list is passed on
// to smtp.SendMail and should match the one in auth.
func (m *Mail) SendMail(adr string, auth smtp.Auth) (err error) {
	var msg []byte
	if msg, err = m.Bytes(); err != nil {
		return
	}
	return smtp.SendMail(adr, auth, m.Recv["From"][0].Address, m.Recipients(), msg)
}

// Formats the mail obj for using a HTML body and returns a buffer that you can
// render your Template to. You must call either HTMLBody or PlainTextBody.
// If you call both, only your last call will be respected.
func (m *Mail) HTMLBody() io.Writer {
	if m.Body == nil {
		m.Body = bytes.NewBuffer(nil)
	}
	m.bodyHeader = "Content-Type: text/html; charset=utf-8\r\n"
	return m.body
}

// Formats the mail obj for using a plaintext body and returns a buffer that you
// can render your Template to. You must call either HTMLBody or PlainTextBody.
// If you call both, only your last call will be respected.
func (m *Mail) PlainTextBody() io.Writer {
	if m.Body == nil {
		m.Body = bytes.NewBuffer(nil)
	}
	m.bodyHeader = "Content-Type: text/plain; charset=utf-8\r\n"
	return m.body
}

// Returns the fully formatted complete message as a slice of bytes.
// Triggers formatting.
func (m *Mail) Bytes() (b []byte, err error) {
	if _, err = m.writeParts(); err != nil {
		return
	}
	b = m.out.Bytes()
	return
}

// Writes the fully formatted complete message to the given writer.
// Triggers formatting.
func (m *Mail) WriteTo(w io.Writer) (n int, err error) {
	if n, err = m.writeParts(); err != nil {
		return
	}
	return w.Write(m.out.Bytes())
}

func (m *Mail) initializeMPHeader() (n int, err error) {
	tmp := make([]byte, 30)
	m.boundary = make([]byte, 64)
	if n, err = io.ReadFull(rand.Reader, tmp); err != nil {
		return
	}
	n = hex.Encode(m.boundary[2:62], tmp)
	m.boundary[0] = 45
	m.boundary[1] = 45
	m.boundary[62] = 13
	m.boundary[63] = 10
	m.mimeHeader = fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\r\n", m.boundary[2:62])
	return
}

func (m *Mail) writeAdrList(header string, adrlist []mail.Address) (n int, err error) {
	stringlist := make([]string, len(adrlist))
	for i, adr := range adrlist {
		stringlist[i] = fmt.Sprintf("%s", &adr)
	}
	if n, err = fmt.Fprintf(m.out, "%s: %s\r\n", header, strings.Join(stringlist, ", ")); err != nil {
		return
	}
	return
}

func (m *Mail) writeHeader() (n int, err error) {
	for Field, address_list := range m.Recv {
		if address_list != nil {
			if n, err = m.writeAdrList(Field, address_list); err != nil {
				return
			}
		}
	}

	if n, err = fmt.Fprintf(m.out, "Subject: %s\r\n", m.Subject); err != nil {
		return
	}

	if n, err = fmt.Fprint(m.out, "MIME-Version: 1.0\r\n"); err != nil {
		return
	}

	if m.Attachments == nil {
		if n, err = fmt.Fprint(m.out, m.bodyHeader, "\r\n"); err != nil {
			return
		}
	} else {
		if n, err = fmt.Fprint(m.out, m.mimeHeader, "\r\n"); err != nil {
			return
		}
	}

	return
}

func (m *Mail) writeBody() (n int, err error) {
	if m.Attachments != nil {
		if n, err = m.out.Write(m.boundary); err != nil {
			return
		}
		if n, err = m.out.WriteString(m.bodyHeader + "\r\n"); err != nil {
			return
		}
	}
	if n, err = m.out.Write(m.body.Bytes()); err != nil {
		return
	}
	if m.Attachments != nil {
		if n, err = m.out.Write(m.boundary[:62]); err != nil {
			return
		}
	}
	return
}

func (m *Mail) writeMIMEFileHeader(buf *bytes.Buffer, filename string) (err error) {
	if _, err = buf.WriteString("Content-Type: application/octet-stream\r\n"); err != nil {
		return
	}
	if _, err = buf.WriteString("Content-Transfer-Encoding: base64\r\n"); err != nil {
		return
	}
	if _, err = buf.WriteString("Content-Disposition: attachment; filename=\"" + filepath.Base(filename) + "\"\r\n\r\n"); err != nil {
		return
	}
	return

}

func (m *Mail) writeAttachment(buf *bytes.Buffer, filename string) (n int, err error) {
	if err = m.writeMIMEFileHeader(buf, filename); err != nil {
		return
	}

	var fObj *os.File
	if fObj, err = os.Open(filename); err != nil {
		return
	}
	enc := base64.NewEncoder(base64.StdEncoding, buf)

	if _, err = io.Copy(enc, fObj); err != nil {
		return
	}
	enc.Close()
	if n, err = buf.Write([]byte{13, 10}); err != nil {
		return
	}
	return
}

func (m *Mail) writeAttachments() (n int, err error) {
	if n, err = m.out.WriteString("\r\n"); err != nil {
		return
	}

	for i, filename := range m.Attachments {
		buffer := bytes.NewBuffer(nil)

		if n, err = m.writeAttachment(buffer, filename); err != nil {
			return
		}

		if i != len(m.Attachments)-1 {
			if n, err = buffer.Write(m.boundary); err != nil {
				return
			}
		} else {
			if n, err = buffer.Write(m.boundary[:62]); err != nil {
				return
			}
		}

		if _, err = io.Copy(m.out, buffer); err != nil {
			return
		}
	}

	if n, err = m.out.Write([]byte{45, 45}); err != nil {
		return
	}
	return
}

func (m *Mail) writeParts() (n int, err error) {
	if m.out == nil {
		m.out = bytes.NewBuffer(nil)
	}
	if n, err = m.initializeMPHeader(); err != nil {
		return
	}

	if n, err = m.writeHeader(); err != nil {
		return
	}

	if m.body != nil {
		if n, err = m.writeBody(); err != nil {
			return
		}
	}

	if m.Attachments != nil {
		if n, err = m.writeAttachments(); err != nil {
			return
		}
	}
	return
}
