package MIMEMail

import (
	"html/template"
	"net/mail"
	"net/smtp"
)

func Example() {
	// create the Mail
	m := NewMail()

	// add Mail addresses to Address fields.
	m.From("foobar", "foobar@example.com")

	// you also can add mail.Address structs
	addr := mail.Address{Name: "baz", Address: "baz@example.com"}
	m.ToAddr(addr)

	// set the subject
	m.Subject = "你好 ma"

	tmpl, err := template.ParseFiles("mailBody.html")
	if err != nil {
		return
	}

	// render your template into the mail body
	if err = tmpl.ExecuteTemplate(m.HTMLBody(), "body", nil); err != nil {
		return
	}

	auth := smtp.PlainAuth("", "foobar@example.com", "foobars password", "mail.example.com")

	// directly send the mail via smtp.SendMail (uses StartTLS if available on the server).
	if err := m.SendMail("mail.example.com:25", auth); err != nil {
		return
	}

	// alternatively, send the mail via TLSClient
	cnf := &Config{
		Host: "mail.example.com",
		Port: "465",
		Auth: auth,
	}
	c, err := TLSClient(cnf)
	if err != nil {
		return
	}

	if err := c.Send(m); err != nil {
		return
	}
}
