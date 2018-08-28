package MIMEMail

import (
	//"bytes"
	"bytes"
	"html/template"
	"net/mail"

	//"os"
	"testing"
)

func MessageFactory() *Mail {
	m := NewMail()
	m.AddPerson("From", "你好 ma", "foobar@example.com")
	m.AddAddress("To", mail.Address{Name: "Änja Süße", Address: "blabla@example.com"})
	m.AddAddress("To", mail.Address{Name: "xiao mao", Address: "xiao_mao@example.com"})
	m.Subject = "你好 Änja"
	return m
}

func Test_Message_writeHeader(t *testing.T) {
	m := MessageFactory()

	b, err := m.Bytes()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("\n", string(b))
}

func Test_Message_Body(t *testing.T) {
	m := MessageFactory()
	tmpl, err := template.New("body").Parse(htmlBody)
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
	m := MessageFactory()

	if err := m.AddReader("short_attachment.txt", bytes.NewBuffer([]byte(shortAttachment))); err != nil {
		t.Fatalf("opening attachment: %s", err)
	}

	tmpl, err := template.New("body").Parse(htmlBody)
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

const (
	htmlBody = `{{ define "body" }}<html>
	<body>
		<h1 id="bla">你好 world!</h1>
		<p>I heard you like MIME Mails, so I put</p>
		<ul>
	        <li>MIMEHeader in your Mailbody</li>
		</ul>
		<p>So that you can</p>
		<ul>
			<li>send MIMEparts while you Multipart</li>
		</ul>
		<p>you have been pimped!</p>
	</body>
	</html>
	{{ end }}`

	shortAttachment = `I'm a short attachment!
`
)
