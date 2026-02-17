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
		// Warnai Key (sebelum titik dua :)
		if strings.Contains(line, ":") {
			parts := strings.SplitN(line, ":", 2)
			key := parts[0]
			val := parts[1]

			// Key warna Cyan/Biru Muda
			key = strings.Replace(key, "\"", colorCyan+"\"", 1)
			key = strings.Replace(key, "\"", "\""+colorReset, 1) // Tutup quote

			// Warnai Value berdasarkan tipe (sederhana)
			if strings.Contains(val, "\"") {
				// String -> Kuning
				val = strings.Replace(val, "\"", colorYellow+"\"", 1)
				val = strings.Replace(val, "\"", "\""+colorReset, -1) // Tutup quote terakhir
			} else if strings.Contains(val, "true") || strings.Contains(val, "false") {
				// Boolean -> Merah
				val = colorRed + val + colorReset
			} else {
				// Angka/Null -> Hijau
				val = colorGreen + val + colorReset
			}

			line = key + ":" + val
		}
		result = append(result, line)
	}
	return strings.Join(result, "\n")
}
