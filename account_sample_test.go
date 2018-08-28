package MIMEMail

// import "net/smtp"
//
// func init() {
// 	testConfig = &testConf{
// 		Config: &Config{
// 			Host: "mail.example.com",
// 			Port: SMTPTLS,
// 			Auth: smtp.PlainAuth("", "sender@example.com", "mailserverpassword", "mail.example.com"),
// 		},
// 		sender: sender{
// 			address: "sender@example.com",
// 			pass:    "your pgp key password",
// 			key: `
// -----BEGIN PGP PRIVATE KEY BLOCK-----
//
// tons of ASCII armored gibberish here!
// get this data with:
// $ gpg --armor --export-secret-keys <your sending mail address>
// -----END PGP PRIVATE KEY BLOCK-----`,
// 		},
// 		receiver: receiver{
// 			address: "receiver@example.com",
// 			key: `
// -----BEGIN PGP PUBLIC KEY BLOCK-----
//
// lots of base64 encoded data here...
// get this data with:
// $ gpg --armor --export <your receiving mail address>
// -----END PGP PUBLIC KEY BLOCK-----
// `,
// 		},
// 	}
// }
