package MIMEMail

import (
	"bytes"
	"crypto"
	"strings"
	"testing"

	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/packet"
)

func TestPGPRead(t *testing.T) {
	cnf := getTestConfig(t)
	e, err := CreateEntity(cnf.receiver.address, strings.NewReader(cnf.receiver.key))
	if err != nil {
		t.Fatal(err)
	}

	out := bytes.NewBuffer(nil)
	plain, err := Encrypt(out, []*openpgp.Entity{e}, nil, nil, &packet.Config{
		DefaultHash: crypto.SHA256,
	})
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
	recv, err := CreateEntity(cnf.receiver.address, strings.NewReader(cnf.receiver.key))
	if err != nil {
		t.Fatal(err)
	}

	signer, err := CreateSigningEntity("", strings.NewReader(cnf.sender.key), cnf.sender.pass)
	if err != nil {
		t.Fatal(err)
	}

	c, err := TLSClient(cnf.Config)
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

	if err := m.WriteEncrypted(out, []*openpgp.Entity{recv}, signer, nil, nil); err != nil {
		t.Fatal(err)
	}
}

func TestPGPEncrypt(t *testing.T) {
	cnf := getTestConfig(t)

	m := getTestMail(t, cnf)

	recv, err := CreateEntity(m.Addresses[AddrTo][0].Address, strings.NewReader(cnf.receiver.key))
	if err != nil {
		t.Fatal(err)
	}

	signer, err := CreateSigningEntity("", strings.NewReader(cnf.sender.key), cnf.sender.pass)
	if err != nil {
		t.Fatal(err)
	}

	out := bytes.NewBuffer(nil)

	if err := m.WriteEncrypted(out, []*openpgp.Entity{recv}, signer, nil, nil); err != nil {
		t.Fatal(err)
	}

	t.Log(out.String())
}
