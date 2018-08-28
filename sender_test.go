package MIMEMail

import (
	"bytes"
	"html/template"
	"testing"
)

func getTestMail(t *testing.T, cnf *testConf) *Mail {
	m := NewMail()
	m.From("MIMEMail test client", cnf.sender.address)
	m.To("MIMEMail recipient", cnf.receiver.address)

	m.Subject = "Test mail from MIMEMail"

	tmpl, err := template.New("body").Parse(htmlBody)
	if err != nil {
		t.Fatal(err)
	}
	if err := tmpl.ExecuteTemplate(m.HTMLBody(), "body", nil); err != nil {
		t.Fatal(err)
	}

	if err := m.AddReader("short_attachment.txt", bytes.NewBuffer([]byte(shortAttachment))); err != nil {
		t.Fatalf("opening attachment: %s", err)
	}
	return m
}

func TestTLSClientConnect(t *testing.T) {
	s, err := TLSClient(getTestConfig(t).Config)
	if err != nil {
		t.Fatal(err)
	}

	if err := s.Quit(); err != nil {
		t.Fatal(err)
	}
	t.Log(s)
}

func TestTLSClientSend(t *testing.T) {
	cnf := getTestConfig(t)

	m := getTestMail(t, cnf)

	c, err := TLSClient(cnf.Config)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Quit()

	if err := c.Send(m); err != nil {
		t.Fatal(err)
	}
}
