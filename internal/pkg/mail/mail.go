package mail

import (
    "gopkg.in/gomail.v2"
    "log"
)

func SendEmail(toEmail string, subject string, body string) error {
    m := gomail.NewMessage()
	m.SetHeader("From", "zharifzuqi@gmail.com")
    m.SetHeader("To", toEmail)
    m.SetHeader("Subject", subject)
    m.SetBody("text/html", body)

    d := gomail.NewDialer("smtp.gmail.com", 587, "zharifzuqi@gmail.com", "shma ueqg pidf zknk")

    if err := d.DialAndSend(m); err != nil {
        log.Fatalf("Failed to send email: %v", err)
		return err
    }

    log.Println("Email sent successfully!")
	return nil
}