package mail

import (
	"log"
	"strconv"

	"github.com/xcurvnubaim/njajal-gin-golang/internal/configs"
	"gopkg.in/gomail.v2"
)

func SendEmail(toEmail string, subject string, body string) error {
    m := gomail.NewMessage()
	m.SetHeader("From", configs.Config.EMAIL_FROM)
    m.SetHeader("To", toEmail)
    m.SetHeader("Subject", subject)
    m.SetBody("text/html", body)
	port, _ := strconv.Atoi(configs.Config.SMTP_PORT)
    d := gomail.NewDialer(configs.Config.SMTP_HOST, port, configs.Config.SMTP_USER, configs.Config.SMTP_PASSWORD)

    if err := d.DialAndSend(m); err != nil {
        log.Fatalf("Failed to send email: %v", err)
		return err
    }

    log.Println("Email sent successfully!")
	return nil
}