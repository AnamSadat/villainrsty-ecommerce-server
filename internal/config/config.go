package config

import (
	"log"
	"os"
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
}

func MustLoad() Config {
	addr := getEnv("APP_ADDR", ":8080")
	dbURL := mustEnv("DATABASE_URL")
	secret := mustEnv("COOKIE_SERCET")
	accessTTL := mustDuration("ACCESS_TTL", 15*time.Minute)
	refreshTTL := mustDuration("REFRESH_TTL", 7*24*time.Hour)

	return Config{
		Addr:         addr,
		DatabaseUrl:  dbURL,
		CookieSecret: secret,
		AccessTTL:    accessTTL,
		RefreshTTL:   refreshTTL,
		CookieDomain: getEnv("COOKIEE_DOMAIN", "localhost"),
		CookieSecure: getEnv("COOKIE_SECURE", "falase") == "true",
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
