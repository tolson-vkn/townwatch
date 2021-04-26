package receivers

import (
	"errors"
	"fmt"
	smtplib "net/smtp"
	"strconv"

	"github.com/sirupsen/logrus"
)

type SMTP struct {
	senderEmail    string
	recipientEmail string
	url            string
	plainAuth      smtplib.Auth
}

func (smtp *SMTP) Notify(title, message string) error {
	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	emailMsg := []byte(
		fmt.Sprintf(
			"To: %s\r\n"+
				"Subject: %s\r\n"+
				"\r\n"+
				"%s\r\n",
			smtp.recipientEmail,
			title,
			message,
		),
	)

	err := smtplib.SendMail(
		smtp.url,
		smtp.plainAuth,
		smtp.senderEmail,
		[]string{smtp.recipientEmail},
		[]byte(emailMsg),
	)
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Infof("Sent mail to [%s].\n", smtp.recipientEmail)
	return nil
}

//func createSMTP(config map[string]interface{}) (SMTP, error) {
func createSMTP(config map[string]interface{}) (*SMTP, error) {
	var exist bool
	var sp int

	// Probably should type check these but it should always be a string? No?
	// and probably could just use some struct syntax magic but whatever.
	if _, exist = config["sender_email"].(string); exist == false {
		return nil, errors.New("Missing [sender_email] in config.")
	}

	if _, exist = config["recipient_email"].(string); exist == false {
		return nil, errors.New("Missing [recipient_email] in config.")
	}

	if _, exist = config["account_email"].(string); exist == false {
		//return nil, errors.New("Missing account email in config")
		return nil, errors.New("Missing [account_email] in config.")
	}

	if _, exist = config["password"].(string); exist == false {
		return nil, errors.New("Missing [password] in config.")
	}

	if _, exist = config["smtp_server"].(string); exist == false {
		return nil, errors.New("Missing [smtp_server] in config.")
	}

	// Port might be given as string or int so do some special things.
	if _, exist = config["smtp_port"]; exist == false {
		return nil, errors.New("Missing [smtp_port] in config.")
	} else {
		switch port := config["smtp_port"].(type) {
		case int:
			sp = config["smtp_port"].(int)
		case string:
			var err error
			sp, err = strconv.Atoi(port)
			if err != nil {
				return nil, errors.New("Couldn't parse port from [smtp_port]")
			}
		default:
			return nil, errors.New("Cannot parse type of [smtp_server].")
		}
	}

	auth := smtplib.PlainAuth(
		"",
		config["account_email"].(string),
		config["password"].(string),
		config["smtp_server"].(string),
	)

	// Build smtp string
	smtpUrl := fmt.Sprintf("%s:%d", config["smtp_server"], sp)

	smtpReceiver := &SMTP{
		config["sender_email"].(string),
		config["recipient_email"].(string),
		smtpUrl,
		auth,
	}

	return smtpReceiver, nil
}
