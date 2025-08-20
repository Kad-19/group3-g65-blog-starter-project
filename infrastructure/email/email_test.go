package email

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/gomail.v2"
)

type mockDialer struct {
	dialAndSendFunc func(m ...*gomail.Message) error
}

func (m *mockDialer) DialAndSend(msg ...*gomail.Message) error {
	return m.dialAndSendFunc(msg...)
}

func TestEmailService(t *testing.T) {
	fromEmail := "test@example.com"

	t.Run("SendActivationEmail success", func(t *testing.T) {
		dialer := &mockDialer{
			dialAndSendFunc: func(m ...*gomail.Message) error {
				assert.Equal(t, "test@example.com", m[0].GetHeader("From")[0])
				assert.Equal(t, "recipient@example.com", m[0].GetHeader("To")[0])
				assert.Equal(t, "Activate Your Account", m[0].GetHeader("Subject")[0])
				return nil
			},
		}
		service := &EmailService{dialer: dialer, fromEmail: fromEmail}
		err := service.SendActivationEmail("recipient@example.com", "http://example.com/activate")
		assert.NoError(t, err)
	})

	t.Run("SendActivationEmail failure", func(t *testing.T) {
		dialer := &mockDialer{
			dialAndSendFunc: func(m ...*gomail.Message) error {
				return errors.New("dial error")
			},
		}
		service := &EmailService{dialer: dialer, fromEmail: fromEmail}
		err := service.SendActivationEmail("recipient@example.com", "http://example.com/activate")
		assert.Error(t, err)
	})

	t.Run("SendPasswordResetEmail success", func(t *testing.T) {
		dialer := &mockDialer{
			dialAndSendFunc: func(m ...*gomail.Message) error {
				assert.Equal(t, "test@example.com", m[0].GetHeader("From")[0])
				assert.Equal(t, "recipient@example.com", m[0].GetHeader("To")[0])
				assert.Equal(t, "Password Reset Request", m[0].GetHeader("Subject")[0])
				return nil
			},
		}
		service := &EmailService{dialer: dialer, fromEmail: fromEmail}
		err := service.SendPasswordResetEmail("recipient@example.com", "reset-token")
		assert.NoError(t, err)
	})

	t.Run("SendPasswordResetEmail failure", func(t *testing.T) {
		dialer := &mockDialer{
			dialAndSendFunc: func(m ...*gomail.Message) error {
				return errors.New("dial error")
			},
		}
		service := &EmailService{dialer: dialer, fromEmail: fromEmail}
		err := service.SendPasswordResetEmail("recipient@example.com", "reset-token")
		assert.Error(t, err)
	})
}
