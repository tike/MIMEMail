package MIMEMail

import (
	"fmt"
	"html/template"
	"net/mail"
	"net/smtp"

	"github.com/tike/MIMEMail/templated"
)

func Example_mime() {
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
		fmt.Println(err)
	}

	// render your template into the mail body
	if err = tmpl.ExecuteTemplate(m.HTMLBody(), "body", nil); err != nil {
		fmt.Println(err)
	}

	auth := smtp.PlainAuth("", "foobar@example.com", "foobars password", "mail.example.com")

	// directly send the mail via smtp.SendMail (uses StartTLS if available on the server).
	if err := m.SendMail("mail.example.com:25", auth); err != nil {
		return
	}

	// alternatively, send the mail via TLSClient
	cnf := &Account{
		Name: "foo bar",

		Address: "mail@example.com",
		Pass:    "foobars password",

		Server: &Server{
			Host: "mail.example.com",
			Port: "465",
		},
	}

	c, err := TLSClient(cnf)
	if err != nil {
		fmt.Println(err)
	}

	if err := c.Send(m); err != nil {
		fmt.Println(err)
	}
}

func Example_pgp() {
	sender := &Account{
		Name:    "Mr. Sender",
		Address: "sender@example.com",
		Pass:    "sender's mail account password",

		Key: &PGP{
			File: "sender.asc",
			Pass: "sender's key password",
		},

		Server: &Server{
			Host: "mail.example.com",
			Port: "465",
		},
	}

	recipient := &Account{
		Name:    "Mr. Receiver",
		Address: "receiver@example.com",

		Key: &PGP{
			File: "receiver.asc",
		},
	}

	m := NewMail()

	// setup the mail, use which ever syntax fits you better
	m.ToAddr(recipient.Addr())
	m.From(sender.Name, sender.Address)
	m.Subject = "PGP test mail"

	bodyContent := `
	<html>
		<body>
			<h1>Hello Mr. Receiver!</h1>
		</body>
	</html>`
	if _, err := m.HTMLBody().Write([]byte(bodyContent)); err != nil {
		fmt.Println(err)
	}

	cipherText, err := m.Encrypt(recipient, sender)
	if err != nil {
		fmt.Println(err)
	}

	from, err := m.EffectiveSender()
	if err != nil {
		fmt.Println(err)
	}

	c, err := TLSClient(sender)
	if err != nil {
		fmt.Println(err)
	}

	if err := c.Write(from, m.Recipients(), cipherText); err != nil {
		fmt.Println(err)
	}
}

func Example_templated() {
	sender := &Account{
		Name:    "Mr. Sender",
		Address: "sender@example.com",
		Pass:    "sender's mail account password",

		Key: &PGP{
			File: "sender.asc",
			Pass: "sender's key password",
		},

		Server: &Server{
			Host: "mail.example.com",
			Port: "465",
		},
	}

	cnf := templated.Config{Dir: "templated/example", Lang: "en_US"}
	ctx := map[string]interface{}{
		"Name":    "Mr. Receiver",
		"Company": "MIMEMail",
	}
	m, err := NewTemplated(&cnf, "foo", ctx)
	if err != nil {
		return
	}
	m.To("Mr. Receiver", "receiver@example.com")
	m.From("Mr. Sender", "sender@example.com")

	if err := m.SendMail("mail.example.com", sender.Auth()); err != nil {
		return
	}
}

func ExampleMail_WriteEncrypted() {
	sender := &Account{
		Name:    "Mr. Sender",
		Address: "sender@example.com",
		Pass:    "sender's mail account password",

		Key: &PGP{
			File: "sender.asc",
			Pass: "sender's key password",
		},

		Server: &Server{
			Host: "mail.example.com",
			Port: "465",
		},
	}

	recipient := &Account{
		Name:    "Mr. Receiver",
		Address: "receiver@example.com",

		Key: &PGP{
			File: "receiver.asc",
		},
	}

	m := NewMail()

	// setup the mail, use which ever syntax fits you better
	m.ToAddr(recipient.Addr())
	m.From(sender.Name, sender.Address)
	m.Subject = "PGP test mail"

	bodyContent := `
	<html>
		<body>
			<h1>Hello Mr. Receiver!</h1>
		</body>
	</html>`
	if _, err := m.HTMLBody().Write([]byte(bodyContent)); err != nil {
		fmt.Println(err)
	}

	from, err := m.EffectiveSender()
	if err != nil {
		fmt.Println(err)
	}

	c, err := TLSClient(sender)
	if err != nil {
		fmt.Println(err)
	}

	out, err := c.W(from, m.Recipients())
	if err != nil {
		fmt.Println(err)
	}

	if err := m.WriteEncrypted(out, recipient, sender); err != nil {
		fmt.Println(err)
	}
}
