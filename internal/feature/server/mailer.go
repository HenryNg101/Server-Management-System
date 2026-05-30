package server

import (
	"fmt"
	"net/smtp"
	"strings"
)

type MailingUtility interface {
	Send(to []string, subject, body string) error
}

type mailerUtility struct {
	server   string
	port     string
	username string
	password string
	from     string
}

func NewMailer(server, port, username, password, from string) MailingUtility {
	return &mailerUtility{server, port, username, password, from}
}

func (m *mailerUtility) Send(to []string, subject, body string) error {
	addr := fmt.Sprintf("%s:%s", m.server, m.port)

	auth := smtp.PlainAuth("", m.username, m.password, m.server)

	msg := []byte(
		"From: " + m.from + "\r\n" +
			"To: " + strings.Join(to, ",") + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/html; charset=\"UTF-8\"\r\n\r\n" +
			body,
	)

	return smtp.SendMail(addr, auth, m.from, to, msg)
}
