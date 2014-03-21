package MIMEMail

import (
	//"bytes"
	"html/template"
	"net/mail"
	//"os"
	"testing"
)

func Message_Factory() *Mail {
	m := NewMail()
	m.AddPerson("From", "你好 ma", "foobar@example.com")
	m.AddAddress("To", mail.Address{"Äjna Süße", "blabla@example.com"})
	m.Subject = "你好 Änja"
	return m
}

func Test_Message_writeHeader(t *testing.T) {
	m := Message_Factory()

	b, err := m.Bytes()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("\n", string(b))
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

	b, err := m.Bytes()
	if err != nil {
		t.Fatal(err)
	}

	t.Log("\n", string(b))
}

func Test_Message_Attach(t *testing.T) {
	m := Message_Factory()

	for _, name := range []string{"short_attachment.txt", "short_attachment.txt"} {
		if err := m.AddFile(name); err != nil {
			t.Fatalf("opening attachment: %s", err)
		}
	}

	tmpl, err := template.ParseFiles("mailBody.html")
	if err != nil {
		t.Fatal(err)
	}

	if err = tmpl.ExecuteTemplate(m.HTMLBody(), "body", []string{"foo", "bar", "baz"}); err != nil {
		t.Fatal(err)
	}

	b, err := m.Bytes()
	if err != nil {
		t.Fatal(err)
	}

	t.Log("\n", string(b))
}
