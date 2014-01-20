MIMEMail
===============================================================================
convenient formatting (and sending) of MIME formatted emails
---------------------------------------------------------------------
With MIMEMail you can easyly send MIME encoded emails. It supports HTML and Plaintext
bodies created from go templates (or anything that can write to io.Writer).
Adding attachments is of course also supported.

How to
------
Simply create a new Mail struct with NewMail(), add Recipients (To, Cc, Bcc, etc).
Set the Subject, add Attachments (by filename), get a Writer for the body by
calling HTMLBody() or PlainTextBody() and render your template into it.
Finally call Bytes() to obtain the formatted email or WriteTo() to directly
write it to a Writer or send it directly (via smtp.SendMail) through the
Mail.SendMail() method.

See godoc for further details.
