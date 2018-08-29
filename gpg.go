package MIMEMail

import (
	"fmt"
	"io"
	"net/mail"

	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"golang.org/x/crypto/openpgp/packet"
)

// UnpackKey parses the ASCIIArmored data read from r.
func UnpackKey(r io.Reader) (*packet.PublicKey, error) {
	deArmored, err := armor.Decode(r)
	if err != nil {
		return nil, err
	}

	pack, err := packet.Read(deArmored.Body)
	if err != nil {
		return nil, err
	}

	pubKey, ok := pack.(*packet.PublicKey)
	if !ok {
		return nil, fmt.Errorf("not a public key")
	}
	return pubKey, nil
}

// UnPackPrivateKey parses the ASCIIArmored data read from r as a private key.
func UnPackPrivateKey(r io.Reader) (*packet.PrivateKey, error) {
	deArmored, err := armor.Decode(r)
	if err != nil {
		return nil, err
	}

	pack, err := packet.Read(deArmored.Body)
	if err != nil {
		return nil, err
	}

	pubKey, ok := pack.(*packet.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("not a private key")
	}
	return pubKey, nil
}

// AddrToPGPUserID converts the given mail.Address into a pgp.UserId.
func AddrToPGPUserID(addr mail.Address) *packet.UserId {
	return packet.NewUserId(addr.Name, "", addr.Address)
}

// CreateEntity creates a reciepient entity using the given address and
// ASCII armored key read from key.
func CreateEntity(addr string, key io.Reader) (*openpgp.Entity, error) {
	pubKey, err := UnpackKey(key)
	if err != nil {
		return nil, err
	}
	userID := AddrToPGPUserID(mail.Address{})
	prim := true

	return &openpgp.Entity{
		PrimaryKey: pubKey,
		Identities: map[string]*openpgp.Identity{
			userID.Id: &openpgp.Identity{
				Name:   userID.Id,
				UserId: userID,
				SelfSignature: &packet.Signature{
					IsPrimaryId: &prim,
					PreferredSymmetric: []uint8{
						uint8(packet.CipherAES128),
						uint8(packet.CipherAES256),
						uint8(packet.CipherCAST5),
					},
					PreferredHash: []uint8{8, 10},
					//PreferredCompression: []uint8{0},
				},
			},
		},
	}, nil
}

// CreateSigningEntity creates a signing entity with the given address and ASCII armor'ed
// key in (decrypting it using pass, if it is encrypted).
func CreateSigningEntity(addr string, key io.Reader, pass string) (*openpgp.Entity, error) {
	privKey, err := UnPackPrivateKey(key)
	if err != nil {
		return nil, err
	}
	if privKey.Encrypted {
		if err := privKey.Decrypt([]byte(pass)); err != nil {
			return nil, err
		}
	}

	userID := AddrToPGPUserID(mail.Address{})
	prim := true
	return &openpgp.Entity{
		PrimaryKey: &privKey.PublicKey,
		PrivateKey: privKey,
		Identities: map[string]*openpgp.Identity{
			userID.Id: &openpgp.Identity{
				Name:   userID.Id,
				UserId: userID,
				SelfSignature: &packet.Signature{
					IsPrimaryId: &prim,
					PreferredSymmetric: []uint8{
						uint8(packet.CipherAES128),
						uint8(packet.CipherAES256),
						uint8(packet.CipherCAST5),
					},
					PreferredHash: []uint8{8, 10},
					//PreferredCompression: []uint8{0},
				},
			},
		},
	}, nil
}

// newASCIIArmorer returns an ASCII armor encoding writer for pgp messages.
func newASCIIArmorer(w io.Writer) (io.WriteCloser, error) {
	return armor.Encode(w, "PGP MESSAGE", nil)
}

// Encrypt encrypts (and ASCII armor encodes) the data written to the returned writer. Remember to Close the writer when you are done.
func Encrypt(out io.Writer, to []*openpgp.Entity, signer *openpgp.Entity, fileHints *openpgp.FileHints, cnf *packet.Config) (io.WriteCloser, error) {
	arm, err := newASCIIArmorer(out)
	if err != nil {
		return nil, err
	}

	plain, err := openpgp.Encrypt(arm, to, signer, fileHints, cnf)
	if err != nil {
		return nil, err
	}

	return closeWrapper{WriteCloser: plain, close: arm}, nil
}

// Sign signs (and ASCII armor encodes) the data written to the returned writer
// Remember to Close the writer when you are done.
func Sign(out io.Writer, signer *Account) (io.WriteCloser, error) {
	if signer == nil {
		return nil, fmt.Errorf("signer cannot be nil!")
	}

	sign, err := CreateSigningEntity(signer)
	if err != nil {
		return nil, err
	}
	config := signer.Key.Config

	arm, err := newASCIIArmorer(out)
	if err != nil {
		return nil, err
	}

	plain, err := openpgp.Sign(arm, sign, nil, config)
	if err != nil {
		return nil, err
	}

	return closeWrapper{WriteCloser: plain, close: arm}, nil
}

// closeWrapper works arround a quirk in the openpgp implementation.
// The implementation wraps it's writer in a NoOPCloser to guard against closing
// before all data has been written if signer is not nil.
type closeWrapper struct {
	io.WriteCloser
	close io.WriteCloser
}

// Close implements io.Closer
func (c closeWrapper) Close() error {
	if err := c.WriteCloser.Close(); err != nil {
		return err
	}
	return c.close.Close()
}
