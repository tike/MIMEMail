MIMEMail
===============================================================================
convenient formatting of MIME formatted emails
---------------------------------------------------------------------
With MIMEMail you can easyly send MIME encoded emails. It supports HTML and Plaintext
bodies created from go templates (or anything that can write to io.Writer).
Adding attachments is of course supported, too.

How to
------
Simply create a new Mail struct with
* NewMail()
* Add Recipients (To, Cc, Bcc, etc).
* Set the Subject
* Add Attachments (by filename or as io.Reader)
* Get a Writer for the body by calling HTMLBody() and/or PlainTextBody()
* Render your template into it (optionally using the included i18n capable templating mechanism)

Optionally
* pgp sign and/or encrypt your message (with full PGP MIME 1.0 formatting)

Finally call
* Bytes() to obtain the formatted email OR
* WriteTo() to directly write it to a Writer OR
* send it directly (via smtp.SendMail) through the Mail.SendMail() method OR
* send it using the included TLS capable SMTP Client


See godoc for further details.
