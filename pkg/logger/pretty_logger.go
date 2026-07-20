package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"sync"
)

// ANSI Colors
const (
	colorReset  = "\033[0m"
	colorCyan   = "\033[36m"
	colorYellow = "\033[33m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorBlue   = "\033[34m"
	colorGray   = "\033[90m"
)

// PrettyHandler adalah custom writer untuk slog
type PrettyHandler struct {
	handler slog.Handler
	w       io.Writer
	mu      *sync.Mutex
}

func NewPrettyHandler(w io.Writer, opts *slog.HandlerOptions) *PrettyHandler {
	// Kita gunakan JSONHandler asli sebagai basis, tapi kita akan "bajak" outputnya nanti
	return &PrettyHandler{
		handler: slog.NewJSONHandler(w, opts),
		w:       w,
		mu:      &sync.Mutex{},
	}
}

func (h *PrettyHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h *PrettyHandler) Handle(ctx context.Context, r slog.Record) error {
	// 1. Ambil data JSON asli dari slog (yang satu baris itu)
	buf := bytes.NewBuffer(nil)
	subHandler := slog.NewJSONHandler(buf, nil)
	if err := subHandler.Handle(ctx, r); err != nil {
		return err
	}

	// 2. Unmarshal ke map supaya bisa kita format ulang
	var data map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &data); err != nil {
		return err
	}

	// 3. Marshal Indent (Bikin jadi berbaris-baris)
	prettyJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	// 4. Mewarnai Output JSON
	coloredJSON := colorizeJSON(string(prettyJSON))

	// 5. Print ke terminal (Thread safe)
	h.mu.Lock()
	defer h.mu.Unlock()
	fmt.Fprintln(h.w, coloredJSON)

	return nil
}

func (h *PrettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &PrettyHandler{handler: h.handler.WithAttrs(attrs), w: h.w, mu: h.mu}
}

func (h *PrettyHandler) WithGroup(name string) slog.Handler {
	return &PrettyHandler{handler: h.handler.WithGroup(name), w: h.w, mu: h.mu}
}

// Fungsi sederhana untuk mewarnai string JSON
func colorizeJSON(jsonStr string) string {
	lines := strings.Split(jsonStr, "\n")
	var result []string

	for _, line := range lines {
		// Skip empty lines dan bracket lines
		trimmed := strings.TrimSpace(line)
		if trimmed == "{" || trimmed == "}" || trimmed == "" {
			result = append(result, line)
			continue
		}

		// Warnai Key dan Value
		if strings.Contains(line, ":") {
			// Cari posisi key (antara { dan :)
			colonIdx := strings.Index(line, ":")
			if colonIdx == -1 {
				result = append(result, line)
				continue
			}

			keyPart := line[:colonIdx]   // Bagian sebelum :
			valPart := line[colonIdx+1:] // Bagian setelah :

			// Warnai key (cyan)
			keyPart = colorizeKey(keyPart)

			// Warnai value berdasarkan tipe
			valPart = colorizeValue(valPart)

			line = keyPart + ":" + valPart
		}
		result = append(result, line)
	}
	return strings.Join(result, "\n")
}

// colorizeKey mewarnai key JSON menjadi cyan
func colorizeKey(keyPart string) string {
	// Find opening quote
	openIdx := strings.Index(keyPart, "\"")
	if openIdx == -1 {
		return keyPart
	}

	// Find closing quote
	closeIdx := strings.LastIndex(keyPart, "\"")
	if closeIdx == -1 || closeIdx <= openIdx {
		return keyPart
	}

	before := keyPart[:openIdx]
	key := keyPart[openIdx : closeIdx+1]
	after := keyPart[closeIdx+1:]

	return before + colorCyan + key + colorReset + after
}

// colorizeValue mewarnai value JSON sesuai tipenya
func colorizeValue(valPart string) string {
	valPart = strings.TrimSpace(valPart)

	// Remove trailing comma if exists
	hasComma := false
	if strings.HasSuffix(valPart, ",") {
		hasComma = true
		valPart = valPart[:len(valPart)-1]
	}

	var colored string

	// Check tipe value
	if valPart == "null" {
		colored = colorGreen + valPart + colorReset
	} else if valPart == "true" || valPart == "false" {
		colored = colorRed + valPart + colorReset
	} else if strings.HasPrefix(valPart, "\"") && strings.HasSuffix(valPart, "\"") {
		// String -> Kuning
		colored = colorYellow + valPart + colorReset
	} else if _, err := parseNumber(valPart); err == nil {
		// Number -> Hijau
		colored = colorGreen + valPart + colorReset
	} else {
		// Unknown -> biarkan aja
		colored = valPart
	}

	if hasComma {
		colored += ","
	}

	return colored
}

// parseNumber helper untuk check apakah string adalah number
func parseNumber(s string) (float64, error) {
	var num float64
	_, err := fmt.Sscanf(s, "%f", &num)
	return num, err
}
