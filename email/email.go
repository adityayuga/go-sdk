package email

import (
	"context"
	"net/smtp"
	"strings"

	"github.com/adityayuga/go-sdk/errors"
)

type (
	Mailer interface {
		Send(ctx context.Context, to, subject, body string) error
	}

	smtpModule struct {
		config Config
	}

	Config struct {
		From     string `json:"from"`
		FromName string `json:"from_name"`
		Password string `json:"password"`
		Host     string `json:"host"`
		Port     string `json:"port"`
	}
)

func New(ctx context.Context, cfg Config) (Mailer, error) {
	return smtpModule{
		config: cfg,
	}, nil
}

func (m smtpModule) Send(ctx context.Context, to, subject, body string) error {
	fromName := m.config.FromName
	if fromName == "" {
		fromName = m.config.From
	}
	if strings.Contains(fromName, "@") {
		fromName = strings.Split(fromName, "@")[0]
	}

	// Setup email message
	msg := []byte(`From: "` + fromName + `" <` + m.config.From + ">\r\n" +
		"To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n\r\n" +
		body)

	// Send the email
	auth := smtp.PlainAuth("", m.config.From, m.config.Password, m.config.Host)
	err := smtp.SendMail(m.config.Host+":"+m.config.Port, auth, m.config.From, []string{to}, msg)
	if err != nil {
		return errors.New(err.Error())
	}

	return nil
}
