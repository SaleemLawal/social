package mailer

import "embed"

const (
	FromName               = "Social"
	MaxRetries             = 3
	UserInvitationTemplate = "user_invitation.tmpl"
)

//go:embed templates
var Fs embed.FS

type Email struct {
	Username      string
	ToEmail       string
	ActivationURL string
}

type Client interface {
	Send(templateFile string, email *Email, isSandbox bool) error
}
