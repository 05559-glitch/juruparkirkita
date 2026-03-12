package util

import (
	"log"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

type SMTPClient struct {
	From   string
	Dialer *gomail.Dialer
}

func NewSMTP() *SMTPClient {
	host := os.Getenv("EMAIL_HOST")
	sender := os.Getenv("EMAIL_SENDER")
	password := os.Getenv("EMAIL_PASSWORD")
	portStr := os.Getenv("EMAIL_PORT")

	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("Format port email tidak valid: %v", err)
	}
	dialer := gomail.NewDialer(host, port, sender, password)

	return &SMTPClient{
		From:   sender,
		Dialer: dialer,
	}
}

func (s *SMTPClient) SendMail(subject, body string, to ...string) error {
	msg := gomail.NewMessage()
	
	msg.SetHeader("From", s.From)
	msg.SetHeader("To", to...)
	msg.SetHeader("Subject", subject)

	msg.SetBody("text/html", body)

	if err := s.Dialer.DialAndSend(msg); err != nil {
		return err
	}

	return nil
}