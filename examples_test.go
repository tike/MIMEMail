package MIMEMail

import (
	"html/template"
	"net/mail"
	"net/smtp"
)

func Example() {
	m := NewMail()

	m.AddPerson("From", "你好 ma", "foobar@example.com")

	as := mail.Address{"Äjna Süße", "blabla@example.com"}
	m.AddAddress("To", as)

	m.Subject = "你好 Änja"

	tmpl, err := template.ParseFiles("mailBody.html")
	if err != nil {
		return
	}

	if err = tmpl.ExecuteTemplate(m.HTMLBody(), "body", nil); err != nil {
		return
	}
	a := smtp.PlainAuth("", "Username", "Password", "mail.example.com")
	m.SendMail("mail.example.com:25", a)
}
