package MIMEMail

import (
	"crypto/tls"
	"io"
	"net"
	"net/mail"
	"net/smtp"
	"os"
	"strings"

	"golang.org/x/crypto/openpgp/packet"
)

// Account holds values for a mail account
type Account struct {
	// Name will be used as the Name component of the "Name <address>" header field
	Name string

	// Address will be used as the Address component of the "Name <address>" header field
	Address string

	// Pass password to login to the server
	Pass string

	// Key holds PGP related values.
	Key *PGP

	// Server infos used to send mail from this account.
	Server *Server
}

// Auth returns a smtp.PlainAuth.
func (a Account) Auth() smtp.Auth {
	return smtp.PlainAuth("", a.Address, a.Pass, a.Server.Addr())
}

// Addr returns a mail.Address
func (a Account) Addr() mail.Address {
	return mail.Address{Name: a.Name, Address: a.Address}
}

// PGP holds the filesystem location of an ASCII Armored PGP private key.
type PGP struct {
	// File holds the filesystem path from which the key should be read.
	// if Key field below is not empty it will take precedence.
	File string

	// Key is mostly for usage in tests, if it is not empty,
	// Open() will return a strings.Reader reading from Key.
	Key string

	// Pass holds the password to decrypt the key (if it is encrypted).
	Pass string

	// Config holds pgp config values. If it is nil sane defaults will be used.
	*packet.Config
}

func (p PGP) Open() (io.Reader, error) {
	if p.Key != "" {
		return strings.NewReader(p.Key), nil
	}

	f, err := os.Open(p.File)
	if err != nil {
		return nil, err
	}
	return f, nil
}

// Server holds connection details for a mail server
type Server struct {
	Host string
	Port string

	// Config holds TLS connection configuration values.
	// if it is nil, sane defaults will be used.
	*tls.Config
}

// Addr return the joint host:port string.
func (s Server) Addr() string {
	return net.JoinHostPort(s.Host, s.Port)
}
