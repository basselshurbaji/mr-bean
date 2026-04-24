package mailer

import "context"

type Email struct {
	To       string
	Subject  string
	Template string
	Data     any
}

type WelcomeData struct {
	FirstName string
}

type Mailer interface {
	Send(ctx context.Context, email Email) error
}