package mailer

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"text/template"
	"time"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGridMailer struct {
	fromEmail string
	apiKey    string
	client    *sendgrid.Client
}

func NewSendGridMailer(fromEmail, apiKey string) *SendGridMailer {
	client := sendgrid.NewSendClient(apiKey)

	return &SendGridMailer{
		fromEmail: fromEmail,
		apiKey:    apiKey,
		client:    client,
	}
}

func (m *SendGridMailer) Send(templateFile string, email *Email, isSandbox bool) error {
	from := mail.NewEmail(FromName, m.fromEmail)
	to := mail.NewEmail(email.Username, email.ToEmail)

	tmpl, err := template.ParseFS(Fs, "templates/"+templateFile)
	if err != nil {
		return err
	}

	subject := new(bytes.Buffer)
	if err := tmpl.ExecuteTemplate(subject, "subject", email); err != nil {
		return err
	}

	body := new(bytes.Buffer)
	if err := tmpl.ExecuteTemplate(body, "body", email); err != nil {
		return err
	}

	message := mail.NewSingleEmail(from, subject.String(), to, "", body.String())

	message.SetMailSettings(&mail.MailSettings{
		SandboxMode: &mail.Setting{
			Enable: &isSandbox,
		},
	})

	for i := range MaxRetries {
		response, err := m.client.Send(message)
		if err != nil {
			log.Printf("Failed to send email to %v, attempt %d of %d", email.ToEmail, i+1, MaxRetries)
			log.Printf("Error: %v", err.Error())

			// Exponential backoff
			time.Sleep(time.Second * time.Duration(math.Pow(2, float64(i))))
			continue
		}

		if response.StatusCode >= 200 && response.StatusCode < 300 {
			log.Printf("Email sent with status code %v", response.StatusCode)
			return nil
		}
	}

	return fmt.Errorf("Failed to send email to %v after %d attempts", email.ToEmail, MaxRetries)
}
