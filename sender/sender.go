package sender

import (
	"log"
	"templtodo3/config"

	mail "github.com/xhit/go-simple-mail/v2"
)

var EmailServer *mail.SMTPServer

func Init() {

	EmailServer = mail.NewSMTPClient()

	// get email configuration details from environment and error if any are missing

	// SMTP Server
	EmailServer.Host = config.AppConfig.EmailHost
	EmailServer.Port = config.AppConfig.EmailPort
	EmailServer.Username = config.AppConfig.EmailUsername
	EmailServer.Password = config.AppConfig.EmailPassword
	EmailServer.Encryption = mail.EncryptionSTARTTLS
}

func SendEmail(to, subject, body string) error {
	email := mail.NewMSG()

	email.SetFrom(config.AppConfig.EmailFrom).
		AddTo(to).
		SetSubject(subject)

	email.SetBody(mail.TextPlain, body)

	client, err := EmailServer.Connect()
	if err != nil {
		log.Default().Printf("failed to connect to email server: %v\n", err)
		return err
	}
	err = email.Send(client)
	if err != nil {
		log.Default().Printf("failed to send email: %v\n", err)
		return err
	}
	return nil
}
