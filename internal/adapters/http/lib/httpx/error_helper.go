package httpx

import (
	"errors"
	"os"
	"regexp"
	"strings"
)

type WrapOpt struct {
	MaxLen          int  // default 160
	HideInProd      bool // default true
	ForceShowInProd bool // kalau true, tetap tampil di prod (hati-hati)
}

func WrapErrDetails(err error, opt WrapOpt) any {
	if err == nil {
		return nil
	}

	// default settings
	if opt.MaxLen <= 0 {
		opt.MaxLen = 160
	}
	// by default: hide details in prod
	if !opt.ForceShowInProd && opt.HideInProd == false {
		// kalau user tidak set apa-apa, tetap hide di prod
		// (opsional: bisa kamu hapus kalau mau)
		opt.HideInProd = true
	}

	if isProd() && opt.HideInProd && !opt.ForceShowInProd {
		return nil // nil => field details akan hilang karena omitempty
	}

	root := rootCause(err)
	s := prettifyErrString(root.Error())
	return truncate(s, opt.MaxLen)
}

func isProd() bool {
	env := strings.ToLower(strings.TrimSpace(os.Getenv("APP_ENV")))
	return env == "production" || env == "prod"
}

func rootCause(err error) error {
	for err != nil {
		u := errors.Unwrap(err)
		if u == nil {
			return err
		}
		err = u
	}
	return nil
}

func prettifyErrString(s string) string {
	s = strings.TrimSpace(s)

	prefixes := []string{
		"unauthorized:",
		"failed to parse token:",
		"token is malformed:",
		"error:",
	}
	low := strings.ToLower(s)
	for _, p := range prefixes {
		if strings.HasPrefix(low, p) {
			s = strings.TrimSpace(s[len(p):])
			low = strings.ToLower(s)
		}
	}

	ws := regexp.MustCompile(`\s+`)
	return ws.ReplaceAllString(s, " ")
}

func truncate(s string, max int) string {
	s = strings.TrimSpace(s)
	if max <= 0 || len(s) <= max {
		return s
	}
	return s[:max] + "..."
}
