package MIMEMail

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/smtp"
)

// various smtp ports as a shorthand
const (
	SMTP         = "25"
	SMTPTLS      = "465"
	SMTPStartTLS = "587"
)

// Client establishes a connection and handles sending the MIME Message(s).
type Client struct {
	*smtp.Client
	cnf *Account
	tls bool
}

// TLSClient establishes a TLSConnection to the Server described by config.
func TLSClient(cnf *Account) (*Client, error) {
	var config *tls.Config
	if cnf.Server.Config == nil {
		config = &tls.Config{ServerName: cnf.Server.Host}
	}
	tlsCon, err := tls.Dial("tcp", cnf.Server.Addr(), config)
	if err != nil {
		return nil, err
	}

	c, err := smtp.NewClient(tlsCon, cnf.Server.Addr())
	if err != nil {
		return nil, err
	}

	return &Client{Client: c, cnf: cnf, tls: true}, nil
}

// PlainClient uses the standard smtp.Dail, so an unencrypted connection will
// be used. It will check wether STARTTLS is supported by the server
// and use it, if available.
func PlainClient(cnf *Account) (*Client, error) {
	c, err := smtp.Dial(cnf.Server.Addr())
	if err != nil {
		return nil, err
	}

	return &Client{Client: c, cnf: cnf}, nil
}

func (c Client) prolog() error {
	if err := c.Hello("localhost"); err != nil {
		c.Quit()
		return err
	}

	if !c.tls {
		if ok, _ := c.Extension("STARTTLS"); ok {
			var config *tls.Config
			if c.cnf.Server.Config == nil {
				config = &tls.Config{ServerName: c.cnf.Server.Host}
			}
			if err := c.StartTLS(config); err != nil {
				return err
			}
		}
	}

	ok, extra := c.Extension("AUTH")
	if !ok {
		c.Quit()
		return fmt.Errorf("no auth: %t %s", ok, extra)
	}

	if err := c.Auth(c.cnf.Auth()); err != nil {
		c.Quit()
		return err
	}

	return nil
}

// W is the equivalent of Writer, but returns a WriteCloser that you can write
// your message to. Remember to close the writer when you are done writing to it.
func (c Client) W(from string, to []string) (io.WriteCloser, error) {
	if err := c.prolog(); err != nil {
		return nil, err
	}

	if err := c.Mail(from); err != nil {
		return nil, err
	}

	for _, addr := range to {
		if err := c.Rcpt(addr); err != nil {
			return nil, err
		}
	}

	return c.Data()
}

// Write sends the message with the given from / to email addresses.
func (c Client) Write(from string, to []string, msg []byte) error {
	w, err := c.W(from, to)
	if err != nil {
		return err
	}
	defer w.Close()

	if _, err := w.Write(msg); err != nil {
		return err
	}

	return nil
}

// Send sends the given Mail.
// If you have the "Sender" field set, it's first entry is used and
// should match the Address in auth, else the first "From" entry
// is used (with the same restrictions). If both are nil,
// a NoSender error is returned. To Send encrypted mails,
// use the (Client.Write / Mail.Encrypt) or (Client.W / Mail.WriteEncrypted)
// pairs.
func (c Client) Send(m *Mail) error {
	efSender, err := m.EffectiveSender()
	if err != nil {
		return err
	}
	recp := m.Recipients()

	b, err := m.Bytes()
	if err != nil {
		return err
	}

	return c.Write(efSender, recp, b)
}

// SendEncrypted sends the given mail to recipient, encrypting it with recipient's Key (which therefore cannot be nil)
// and signing it with sender's Key (which may also not be nil).
func (c Client) SendEncrypted(m *Mail, recipient, sender *Account) error {
	enc, err := m.Encrypt(recipient, sender)
	if err != nil {
		return err
	}

	return c.Write(sender.Address, []string{recipient.Address}, enc)
}
