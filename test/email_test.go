package test

import (
	"testing"

	"github.com/xcurvnubaim/Task-1-IS/internal/pkg/mail"
)

func TestSendEmail(t *testing.T) {
	err := mail.SendEmail("zhariffmarzuqi@gmail.com", "Test Email from Go", "<a href=\"google.com\">This is a test email sent from a Go application!</a>")
	if err != nil {
		t.Fatalf("Failed to send email: %v", err)
	}
}