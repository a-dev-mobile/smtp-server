package utils

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"net/smtp"
	"github.com/a-dev-mobile/smtp-server/internal/config"
	"strings"

	"golang.org/x/exp/slog"
)

// EmailProvider defines the types of email providers.
type EmailProvider int


// EmailConfig stores the configuration for an email provider.
type EmailConfig struct {
	Name      string
	SMTPHost  string
	SMTPPort  string
	Login     string
	Password  string
	FromEmail string
	Logger    *slog.Logger
}

// NewEmailConfigs creates a slice of EmailConfig from a slice of SMTPProviderConfig.
func NewEmailConfigs(smtpProviders []config.SMTPProviderConfig, lg *slog.Logger) []*EmailConfig {
	var emailConfigs []*EmailConfig
	for _, provider := range smtpProviders {
		emailConfigs = append(emailConfigs, &EmailConfig{
			Name:      provider.Name,
			SMTPHost:  provider.SMTPHost,
			SMTPPort:  provider.SMTPPort,
			Login:     provider.Login,
			Password:  provider.Password,
			FromEmail: provider.FromEmail,
			Logger:    lg,
		})
	}
	return emailConfigs
}

// SendEmail tries to send an email using the provided configurations.
// It iterates through the configurations and attempts to send the email, logging errors as they occur.
func SendEmail(configs []*EmailConfig, fromName, defaultFromEmail, to, subject, body string) error {
	var errors []string
	for i, config := range configs {
		if err := config.validateConfig(); err != nil {
			errors = append(errors, err.Error())
			continue
		}

		fromEmail := chooseFromEmail(config.FromEmail, defaultFromEmail)
		if err := sendEmailWithConfig(config, fromName, fromEmail, to, subject, body); err != nil {
			logError(config.Logger, i, config, err)
			errors = append(errors, err.Error())
			continue
		}
		logSuccess(config.Logger, config)
		return nil
	}
	return fmt.Errorf("all SMTP configurations failed: %s", strings.Join(errors, "; "))
}

// sendEmailWithConfig sends an email using a specific EmailConfig.
// It establishes a TLS connection and sends the email using the SMTP protocol.
func sendEmailWithConfig(config *EmailConfig, fromName, fromEmail, to, subject, body string) error {
	encodedSubject := encodeSubject(subject)
	message := buildEmailMessage(fromName, fromEmail, to, encodedSubject, body)

	tlsConfig := &tls.Config{ServerName: config.SMTPHost}
	auth := smtp.PlainAuth("", config.Login, config.Password, config.SMTPHost)

	addr := config.SMTPHost + ":" + config.SMTPPort
	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %v", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, config.SMTPHost)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %v", err)
	}
	defer client.Close()

	return send(client, auth, fromEmail, to, message)
}

// validateConfig checks if the EmailConfig is properly set.
func (ec *EmailConfig) validateConfig() error {
	if ec.SMTPHost == "" || ec.SMTPPort == "" || ec.Login == "" || ec.Password == "" {
		return fmt.Errorf("incomplete email config for %s", ec.Name)
	}
	return nil
}

// chooseFromEmail selects the appropriate email address to use as the sender.
func chooseFromEmail(configEmail, defaultEmail string) string {
	if configEmail != "" {
		return configEmail
	}
	return defaultEmail
}

// logError logs an error that occurred during the email sending process.
func logError(lg *slog.Logger, attempt int, config *EmailConfig, err error) {
	if lg != nil {
		lg.Error(
			"Email sending failed",
			"attempt", attempt+1,
			"smtpHost", config.SMTPHost,
			"error", err)
	}
}

// logSuccess logs a successful email sending operation.
func logSuccess(lg *slog.Logger, config *EmailConfig) {
	if lg != nil {
		lg.Info("Email sent successfully", "name smtp", config.Name)
	}
}

// encodeSubject encodes the email subject line using base64 to support UTF-8 characters.
func encodeSubject(subject string) string {
	return "=?UTF-8?B?" + base64.StdEncoding.EncodeToString([]byte(subject)) + "?="
}

// buildEmailMessage constructs the email message to be sent.
func buildEmailMessage(fromName, fromEmail, to, subject, body string) []byte {
	fromHeader := fmt.Sprintf("From: %s<%s>", fromName, fromEmail)
	toHeader := "To: " + to
	subjectHeader := "Subject: " + subject
	contentTypeHeader := "MIME-Version: 1.0\r\n" +
		"Content-Type: text/plain; charset=\"UTF-8\""
	return []byte(fromHeader + "\r\n" + toHeader + "\r\n" + subjectHeader + "\r\n" +
		contentTypeHeader + "\r\n\r\n" + body + "\r\n")
}

// send finalizes the email sending process by setting the sender, recipient, and sending the message.
func send(client *smtp.Client, auth smtp.Auth, fromEmail, to string, message []byte) error {
	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("authentication failed: %v", err)
	}
	if err := client.Mail(fromEmail); err != nil {
		return fmt.Errorf("failed to set sender: %v", err)
	}
	if err := client.Rcpt(to); err != nil {
		return fmt.Errorf("failed to set recipient: %v", err)
	}

	wc, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to create data writer: %v", err)
	}

	_, err = wc.Write(message)
	if err != nil {
		return fmt.Errorf("failed to write message: %v", err)
	}

	err = wc.Close()
	if err != nil {
		return fmt.Errorf("failed to close writer: %v", err)
	}

	return client.Quit()
}
