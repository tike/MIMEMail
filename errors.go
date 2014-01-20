package MIMEMail

import "fmt"

type NoSender int

func (e NoSender) Error() string {
	return fmt.Sprintf("You have neither From nor Sender set on your Mail!")
}
