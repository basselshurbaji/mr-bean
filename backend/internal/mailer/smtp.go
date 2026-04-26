package mailer

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"html/template"
	"net/smtp"
)

//go:embed templates/*.html
var templateFS embed.FS

var templates = template.Must(template.ParseFS(templateFS, "templates/*.html"))

type smtpMailer struct {
	host     string
	port     string
	username string
	password string
	from     string
}

func NewSMTPMailer(host, port, username, password, from string) Mailer {
	return &smtpMailer{
		host:     host,
		port:     port,
		username: username,
		password: password,
		from:     from,
	}
}

// Send implements Mailer.
func (m *smtpMailer) Send(_ context.Context, email Email) error {
	var body bytes.Buffer
	if err := templates.ExecuteTemplate(&body, email.Template, email.Data); err != nil {
		return fmt.Errorf("render template %q: %w", email.Template, err)
	}

	raw := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s",
		m.from, email.To, email.Subject, body.String(),
	)

	auth := smtp.PlainAuth("", m.username, m.password, m.host)
	return smtp.SendMail(m.host+":"+m.port, auth, m.from, []string{email.To}, []byte(raw))
}
