package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Addr         string
	DatabaseUrl  string
	CookieSecret string

	AccessTTL  time.Duration
	RefreshTTL time.Duration

	CookieDomain string
	CookieSecure bool

	SMTPHost      string
	SMTPPort      string
	SMTPUsername  string
	SMTPPassword  string
	SMTPFromEmail string
	SMTPFromName  string

	ResetPasswordTTL time.Duration
	ResetPasswordURL string
	TwoFactorOTPTTL  time.Duration
}

func MustLoad() Config {
	addr := getEnv("APP_ADDR", ":8080")
	dbURL := mustEnv("DATABASE_URL")
	secret := mustEnv("COOKIE_SECRET")
	accessTTL := mustDuration("ACCESS_TTL", 15*time.Minute)
	refreshTTL := mustDuration("REFRESH_TTL", 7*24*time.Hour)
	smtpHost := mustEnv("SMTP_HOST")
	smtpPort := mustEnv("SMTP_PORT")
	smtpUsername := mustEnv("SMTP_USERNAME")
	smtpPassword := mustEnv("SMTP_PASSWORD")
	smtpFromEmail := mustEnv("SMTP_FROM_EMAIL")
	smtpFromName := mustEnv("SMTP_FROM_NAME")
	resetPasswordTTL := mustDuration("RESET_PASSWORD_TTL", 30*time.Minute)
	resetPasswordURL := mustEnv("FRONTEND_RESET_PASSWORD_URL")
	twoFactorOTPTTL := mustDuration("TWO_FACTOR_OTP_TTL", 5*time.Minute)

	return Config{
		Addr:             addr,
		DatabaseUrl:      dbURL,
		CookieSecret:     secret,
		AccessTTL:        accessTTL,
		RefreshTTL:       refreshTTL,
		CookieDomain:     getEnv("COOKIEE_DOMAIN", "localhost"),
		CookieSecure:     getEnv("COOKIE_SECURE", "false") == "true",
		SMTPHost:         smtpHost,
		SMTPPort:         smtpPort,
		SMTPUsername:     smtpUsername,
		SMTPPassword:     smtpPassword,
		SMTPFromEmail:    smtpFromEmail,
		SMTPFromName:     smtpFromName,
		ResetPasswordTTL: resetPasswordTTL,
		ResetPasswordURL: resetPasswordURL,
		TwoFactorOTPTTL:  twoFactorOTPTTL,
	}
}

func getEnv(k, def string) string {
	v := os.Getenv(k)
	if v == "" {
		return def
	}

	return v
}

func mustEnv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("missing env %s", k)
	}

	return v
}

func mustDuration(k string, def time.Duration) time.Duration {
	v := os.Getenv(k)
	if v == "" {
		return def
	}

	d, err := time.ParseDuration(v)
	if err != nil {
		log.Fatalf("invalid duration %s=%q: %v", k, v, err)
	}

	return d
}

func mustBool(k string, def bool) bool {
	v := os.Getenv(k)
	if v == "" {
		return def
	}

	b, err := strconv.ParseBool(v)
	if err != nil {
		log.Fatalf("invalid bool %s=%q: %v", k, v, err)
	}

	return b
}
