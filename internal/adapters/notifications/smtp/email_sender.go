package smtp

import (
	"context"
	"fmt"
	"log"
	"net/smtp"
	"strings"
	"time"

	"villainrsty-ecommerce-server/internal/core/shared/errors"
)

type EmailSender struct {
	host      string
	port      string
	username  string
	password  string
	fromEmail string
	fromName  string
}

func NewEmailSender(host, port, username, password, fromEmail, fromName string) *EmailSender {
	return &EmailSender{
		host:      host,
		port:      port,
		username:  username,
		password:  password,
		fromEmail: fromEmail,
		fromName:  fromName,
	}
}

func (s *EmailSender) SendPasswordReset(ctx context.Context, toEmail, resetLink string) error {
	if toEmail == "" {
		return errors.New(errors.ErrValidation, "email is required")
	}

	if resetLink == "" {
		return errors.New(errors.ErrValidation, "reset link is required")
	}

	subject := "Reset Password"
	plainBody := fmt.Sprintf("klik link berikut untuk reset password Anda: %s\nLink berlaku sementara.", resetLink)
	htmlBody := fmt.Sprintf(`<p>Klik link berikut untuk reset password Anda:</p><p><a href="%s">Reset Password</a></p><p>Link berlaku sementara</p>`, resetLink)

	msg := buildMIMEMessage(s.fromName, s.fromEmail, toEmail, subject, plainBody, htmlBody)
	addr := fmt.Sprintf("%s:%s", s.host, s.port)
	auth := smtp.PlainAuth("", s.username, s.password, s.host)

	start := time.Now()
	deadline, hasDeadline := ctx.Deadline()
	log.Printf("[smtp] start to=%s from=%s user=%s addr=%s host=%s port=%s subject=%q msgBytes=%d hasDeadline=%v deadline=%v",
		toEmail, s.fromEmail, s.username, addr, s.host, s.port, subject, len(msg), hasDeadline, deadline)

	done := make(chan error, 1)
	go func() {
		err := smtp.SendMail(addr, auth, s.fromEmail, []string{toEmail}, []byte(msg))
		done <- err
	}()

	select {
	case <-ctx.Done():
		log.Printf("[smtp] ctx cancelled after=%s err=%v", time.Since(start), ctx.Err())

		return errors.Wrap(errors.ErrInternal, "email send cancelled", ctx.Err())
	case err := <-done:
		if err != nil {
			log.Printf("[smtp] send FAIL after=%s to=%s from=%s addr=%s err=%v", time.Since(start), toEmail, s.fromEmail, addr, err)

			return errors.Wrap(errors.ErrInternal, "failed to send email", err)
		}
		log.Printf("[smtp] send OK after=%s to=%s from=%s addr=%s", time.Since(start), toEmail, s.fromEmail, addr)

		return nil
	case <-time.After(30 * time.Second):
		log.Printf("[smtp] TIMEOUT after=%s to=%s from=%s addr=%s", time.Since(start), toEmail, s.fromEmail, addr)

		return errors.New(errors.ErrInternal, "smtp timeout")
	}
}

func (s *EmailSender) SendLoginOTP(ctx context.Context, toEmail, otpCode string) error {
	if toEmail == "" {
		return errors.New(errors.ErrValidation, "email is required")
	}

	if otpCode == "" {
		return errors.New(errors.ErrValidation, "otp code is required")
	}

	subject := "Kode OTP Login"
	plainBody := fmt.Sprintf("Kode OTP login Anda: %s\nKode berlaku selama beberapa menit.", otpCode)
	htmlBody := fmt.Sprintf(`<p>Kode OTP Login Anda:</p><h2>%s</h2><p>Kode berlaku selama beberapa menit.</p>`, otpCode)

	msg := buildMIMEMessage(s.fromName, s.fromEmail, toEmail, subject, plainBody, htmlBody)
	addr := fmt.Sprintf("%s:%s", s.host, s.port)
	auth := smtp.PlainAuth("", s.username, s.password, s.host)

	done := make(chan error, 1)
	go func() {
		done <- smtp.SendMail(addr, auth, s.fromEmail, []string{toEmail}, []byte(msg))
	}()

	select {
	case <-ctx.Done():
		return errors.Wrap(errors.ErrInternal, "email send cancalled", ctx.Err())
	case err := <-done:
		if err != nil {
			return errors.Wrap(errors.ErrInternal, "failed to send email", err)
		}
		return nil

	case <-time.After(30 * time.Second):
		return errors.New(errors.ErrInternal, "smtp timeout")
	}
}

func buildMIMEMessage(frontName, fromEmail, toEmail, subject, plainBody, htmlBody string) string {
	boundary := "mixed-boundary-reset-password"

	headers := []string{
		fmt.Sprintf("From %s <%s>", frontName, fromEmail),
		fmt.Sprintf("To: %s", toEmail),
		fmt.Sprintf("Subject: %s", subject),
		"MIME-Version: 1.0",
		fmt.Sprintf("Content-Type: multipart/alternative; boundary=%q", boundary),
		"",
	}

	body := []string{
		fmt.Sprintf("--%s", boundary),
		"Content-Type: text/plain; charset=UTF-8",
		"",
		plainBody,
		fmt.Sprintf("--%s", boundary),
		"Content-Type: text/html; charset=UTF-8",
		"",
		htmlBody,
		fmt.Sprintf("--%s--", boundary),
		"",
	}

	return strings.Join(append(headers, body...), "\r\n")
}
