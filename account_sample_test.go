package MIMEMail

// func init() {
// 	testConfig = &testConf{
// 		sender: &Account{
// 			Name:    "Mr. Sender",
// 			Address: "sender@example.com",
// 			Pass:    "your mail server password",
// 			Key: &PGP{
// 				Pass: "my pgp key password",
// 				Key: `
// -----BEGIN PGP PRIVATE KEY BLOCK-----
// tons of ASCII armored gibberish here!
// get this data with:
// $ gpg --armor --export-secret-keys <your sending mail address>
// -----END PGP PRIVATE KEY BLOCK-----
// `
// 			},
// 			Server: &Server{
// 				Host: "mail.example.com",
// 				Port: SMTPTLS,
// 			},
// 		},
// 		receiver: &Account{
// 			Name:    "Mr. Receiver",
// 			Address: "receiver@example.com",
// 			Key: &PGP{
// 				File: "receivers_public_key.asc",
// 				// file contains:`
// 				// -----BEGIN PGP PUBLIC KEY BLOCK-----
// 				//
// 				// lots of base64 encoded data here...
// 				// get this data with:
// 				// $ gpg --armor --export <your receiving mail address>
// 				// -----END PGP PUBLIC KEY BLOCK-----
// 			},
// 		},
// 	}
// }
