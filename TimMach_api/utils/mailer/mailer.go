package mailer

import (
	"errors"
	"fmt"
	"io"

	"gopkg.in/gomail.v2"
)

// Attachment represents a file attached to the email.
type Attachment struct {
	Filename string
	MimeType string
	Content  []byte
}

// Mailer wraps gomail dialer plus metadata about sender and config state.
type Mailer struct {
	dialer  *gomail.Dialer
	from    string
	enabled bool
}

var ErrNotConfigured = errors.New("mailer is not configured")

// New creates a Mailer. If host/user/pass are empty, the mailer is disabled
// and Send will return ErrNotConfigured.
func New(host string, port int, username, password, from string) *Mailer {
	if host == "" || username == "" || password == "" {
		return &Mailer{enabled: false}
	}
	if from == "" {
		from = username
	}
	dialer := gomail.NewDialer(host, port, username, password)
	return &Mailer{
		dialer:  dialer,
		from:    from,
		enabled: true,
	}
}

func (m *Mailer) Enabled() bool {
	return m != nil && m.enabled
}

// Send delivers an email with optional attachments.
func (m *Mailer) Send(to, subject, body string, attachments []Attachment) error {
	if !m.Enabled() {
		return ErrNotConfigured
	}
	if subject == "" {
		subject = "(no subject)"
	}
	msg := gomail.NewMessage()
	msg.SetHeader("From", m.from)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/plain", body)

	for _, att := range attachments {
		name := att.Filename
		if name == "" {
			name = "attachment"
		}
		mime := att.MimeType
		if mime == "" {
			mime = "application/octet-stream"
		}

		msg.Attach(
			name,
			gomail.SetHeader(map[string][]string{
				"Content-Type": {mime},
			}),
			gomail.SetCopyFunc(func(w io.Writer) error {
				_, err := w.Write(att.Content)
				if err != nil {
					return fmt.Errorf("write attachment %s: %w", name, err)
				}
				return nil
			}),
		)
	}

	return m.dialer.DialAndSend(msg)
}
