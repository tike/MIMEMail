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

// Config bundles the relevant bits of information
type Config struct {
	Host string
	Port string
	Auth smtp.Auth

	TLSConf *tls.Config
}

// String implements fmt.Stringer
func (c Config) String() string {
	return c.Net()
}

// Net returns Host:Port
func (c Config) Net() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

// TLSConfig returns the tls.Config (or a sane default if TLSConfig is nil)
func (c *Config) TLSConfig() *tls.Config {
	if c.TLSConf != nil {
		return c.TLSConf
	}
	return &tls.Config{ServerName: c.Host}
}

// Client establisches a TLSConnection and handles sending the MIME Message(s).
type Client struct {
	*smtp.Client
}

// TLSClient establishes a TLSConnection to the Server described by config,
// performs the Hello() and Auth() calls and returns the "ready to send" server.
// If an error occurs, the server connection will be closed and the error returned.
// When you want to run the test, remember to export the necessary env vars to inject
// your config.
func TLSClient(cnf *Config) (*Client, error) {
	tlsCon, err := tls.Dial("tcp", cnf.Net(), cnf.TLSConfig())
	if err != nil {
		return nil, err
	}

	c, err := smtp.NewClient(tlsCon, cnf.Host)
	if err != nil {
		return nil, err
	}

	if err = c.Hello("localhost"); err != nil {
		c.Quit()
		return nil, err
	}

	ok, extra := c.Extension("AUTH")
	if !ok {
		c.Quit()
		return nil, fmt.Errorf("no auth: %t %s", ok, extra)
	}

	if err = c.Auth(cnf.Auth); err != nil {
		c.Quit()
		return nil, err
	}

	return &Client{Client: c}, nil
}

func (c Client) W(from string, to []string) (io.WriteCloser, error) {
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

func (c Client) Write(from string, to []string, data []byte) error {
	w, err := c.W(from, to)
	if err != nil {
		return err
	}
	defer w.Close()

	if _, err := w.Write(data); err != nil {
		return err
	}

	return nil
}

// Send sends the given Mail.
// If you have the "Sender" field set, it's first entry is used and
// should match the Address in auth, else the first "From" entry
// is used (with the same restrictions). If both are nil,
// a NoSender error is returned.
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
