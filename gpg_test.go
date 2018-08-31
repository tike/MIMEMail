package MIMEMail

import (
	"bytes"
	"testing"
)

func TestPGPRead(t *testing.T) {
	cnf := getTestConfig(t)

	out := bytes.NewBuffer(nil)
	plain, err := Encrypt(out, cnf.receiver, nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := plain.Write([]byte("foo bar baz, es pgp'ed so schön!")); err != nil {
		t.Fatal(err)
	}
	plain.Close()

	t.Logf("encrypted:\n%s\n", out.Bytes())
}

func TestPGPSend(t *testing.T) {
	cnf := getTestConfig(t)
	m := getTestMail(t, cnf)

	c, err := TLSClient(cnf.sender)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Quit()

	from, err := m.EffectiveSender()
	if err != nil {
		t.Fatal(err)
	}

	out, err := c.W(from, m.Recipients())
	if err != nil {
		t.Fatal(err)
	}

	if err := m.WriteEncrypted(out, cnf.receiver, cnf.sender); err != nil {
		t.Fatal(err)
	}
}

func TestPGPEncrypt(t *testing.T) {
	cnf := getTestConfig(t)

	m := getTestMail(t, cnf)

	out := bytes.NewBuffer(nil)

	if err := m.WriteEncrypted(out, cnf.receiver, cnf.sender); err != nil {
		t.Fatal(err)
	}

	t.Log(out.String())
}
