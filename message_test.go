package MIMEMail

import (
	"html/template"
	"net/mail"
	"os"
	"testing"
)

func Message_Factory() *Mail {
	m := NewMail()
	m.From = []*mail.Address{TMX.Address}
	m.Bcc = []*mail.Address{GMX.Address}
	m.Subject = "Jo Doag!"
	return m
}

func Test_Message_writeHeader(t *testing.T) {
	m := Message_Factory()
	if n, err := m.WriteTo(os.Stdout); err != nil {
		t.Fatal(n, err)
	}
}

func Test_Message_Body(t *testing.T) {
	m := Message_Factory()
	tmpl, err := template.ParseFiles("mailBody.html")
	if err != nil {
		t.Fatal(err)
	}

	if err = tmpl.ExecuteTemplate(m.HTMLBody(), "body", nil); err != nil {
		t.Fatal(err)
	}

	if n, err := m.WriteTo(os.Stdout); err != nil {
		t.Fatal(n, err)
	}
}

func Test_Message_Attach(t *testing.T) {
	m := Message_Factory()
	m.Attachments = []string{"C:/Users/tike/Pictures/ww.jpg", "C:/Users/tike/Pictures/pim.jpg"}
	tmpl, err := template.ParseFiles("mailBody.html")
	if err != nil {
		t.Fatal(err)
	}

	if err = tmpl.ExecuteTemplate(m.HTMLBody(), "body", []string{"foo", "bar", "baz"}); err != nil {
		t.Fatal(err)
	}
	if n, err := m.WriteTo(os.Stdout); err != nil {
		t.Fatal(n, err)
	}
}

func Test_SendMail(t *testing.T) {
	if err = m.SendMail(ARCOR.HostAdr(), ARCOR.Auth()); err == nil {
		t.Fatal(err)
	}
}