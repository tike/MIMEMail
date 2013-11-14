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
	/* Address Lists for the Mailheader,
	   Fields that are nil, will be ignored. */
	Sender     []*mail.Address
	From       []*mail.Address
	To         []*mail.Address
	Cc         []*mail.Address
	Bcc        []*mail.Address
	ReplyTo    []*mail.Address
	FollowupTo []*mail.Address

	//The subject Line
	Subject string

	//Filenames of Attachments to send along
	Attachments []string

	boundary   []byte
	mimeHeader string
	bodyHeader string

	body *bytes.Buffer
	out  *bytes.Buffer
}

func NewMail() *Mail {
	return &Mail{out: bytes.NewBuffer(nil)}
}

func (m *Mail) Recipients() (to []string) {
	to = make([]string, 0, 10)
	for _, mail := range m.To {
		to = append(to, mail.Address)
	}
	for _, mail := range m.Cc {
		to = append(to, mail.Address)
	}
	for _, mail := range m.Bcc {
		to = append(to, mail.Address)
	}
	return
}

func (m *Mail) SendMail(adr string, a smtp.Auth) (err error) {
	var msg []byte
	if msg, err = m.Bytes(); err != nil {
		return
	}
	return smtp.SendMail(adr, a, m.From[0].Address, m.Recipients(), msg)
}

func (m *Mail) HTMLBody() io.Writer {
	m.body = bytes.NewBuffer(nil)
	m.bodyHeader = "Content-Type: text/html; charset=utf-8\r\n"
	return m.body
}

func (m *Mail) PlainTextBody() io.Writer {
	m.body = bytes.NewBuffer(nil)
	m.bodyHeader = "Content-Type: text/plain; charset=utf-8\r\n"
	return m.body
}

func (m *Mail) Bytes() (b []byte, err error) {
	if _, err = m.writeParts(); err != nil {
		return
	}
	b = m.out.Bytes()
	return
}

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

func (m *Mail) writeAdrList(header string, adrlist []*mail.Address) (n int, err error) {
	stringlist := make([]string, len(adrlist))
	for i, adr := range adrlist {
		stringlist[i] = fmt.Sprintf("%s", adr)
	}
	if n, err = fmt.Fprintf(m.out, "%s: %s\r\n", header, strings.Join(stringlist, ", ")); err != nil {
		return
	}
	return
}

func (m *Mail) writeHeader() (n int, err error) {
	if m.Sender != nil {
		if n, err = m.writeAdrList("Sender: %s\r\n", m.From); err != nil {
			return
		}
	}

	if m.From != nil {
		if n, err = m.writeAdrList("From", m.From); err != nil {
			return
		}
	}

	if m.To != nil {
		if n, err = m.writeAdrList("To", m.To); err != nil {
			return
		}
	}

	if m.Cc != nil {
		if n, err = m.writeAdrList("Cc", m.Cc); err != nil {
			return
		}
	}

	if m.Bcc != nil {
		if n, err = m.writeAdrList("Bcc", m.Bcc); err != nil {
			return
		}
	}

	if m.ReplyTo != nil {
		if n, err = m.writeAdrList("Reply-To", m.ReplyTo); err != nil {
			return
		}
	}

	if m.FollowupTo != nil {
		if n, err = m.writeAdrList("Mail-Followup-To", m.FollowupTo); err != nil {
			return
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
