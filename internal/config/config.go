package config

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Addr         string
	Host         string
	DatabaseUrl  string
	CookieSecret string

	AccessTTL  time.Duration
	RefreshTTL time.Duration

	SMTPHost      string
	SMTPPort      string
	SMTPUsername  string
	SMTPPassword  string
	SMTPFromEmail string
	SMTPFromName  string

	ResetPasswordTTL time.Duration
	ResetPasswordURL string
	TwoFactorOTPTTL  time.Duration

	CasbinModelPath  string
	CasbinPolicyPath string
}

func MustLoad() Config {
	addr := getEnv("APP_ADDR", "8080")
	host := getEnv("APP_HOST", "localhost")
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
	casbinModelPath := mustEnv("CASBIN_MODEL_PATH")
	casbinPolicyPath := mustEnv("CASBIN_POLICY_PATH")

	return Config{
		Addr:             addr,
		Host:             host,
		DatabaseUrl:      dbURL,
		CookieSecret:     secret,
		AccessTTL:        accessTTL,
		RefreshTTL:       refreshTTL,
		SMTPHost:         smtpHost,
		SMTPPort:         smtpPort,
		SMTPUsername:     smtpUsername,
		SMTPPassword:     smtpPassword,
		SMTPFromEmail:    smtpFromEmail,
		SMTPFromName:     smtpFromName,
		ResetPasswordTTL: resetPasswordTTL,
		ResetPasswordURL: resetPasswordURL,
		TwoFactorOTPTTL:  twoFactorOTPTTL,
		CasbinModelPath:  casbinModelPath,
		CasbinPolicyPath: casbinPolicyPath,
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

func splitCSV(v string) []string {
	v = strings.TrimSpace(v)
	if v == "" {
		return nil
	}

	parts := strings.Split(v, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(strings.ToLower(p))
		if p != "" {
			out = append(out, p)
		}
	}

	return out
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
