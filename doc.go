// Package MIMEMail provides convenient formatting (and sending) of MIME formatted emails.
//
// Simply create a new Mail struct with NewMail(), add Recipients (To, Cc, Bcc, etc). Set
// the Subject, add Attachments (by filename), get a Writer for the body
// by calling HTMLBody() or PlainTextBody() and render your template into it.
// Finally call Bytes() to obtain the formatted email or WriteTo() to directly
// write it to a Writer or send it directly (via smtp.SendMail) through the
// Mail.SendMail() method.
//
// You can put your mail templates into files and load / execute them using the
// NewTemplated() constructor (take a look at the templated package to get an
// an example of how to structure/organise your template folder.
//
// Additionally an easy interface for encrypting the message using PGPMime is provided.
package MIMEMail
