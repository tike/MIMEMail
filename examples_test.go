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
	m.AddPerson("From", "foobar", "foobar@example.com")

	// you also can add mail.Address structs
	as := mail.Address{Name: "baz", Address: "baz@example.com"}
	m.AddAddress("To", as)

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

	a := smtp.PlainAuth("", "foobar", "foobars password", "foobar.example.com")

	// directly send the mail.
	if err := m.SendMail("mail.example.com:25", a); err != nil {
		return
	}
}
