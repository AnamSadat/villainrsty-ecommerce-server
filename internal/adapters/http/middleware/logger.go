package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

// ANSI colors
const (
	reset  = "\033[0m"
	green  = "\033[32m"
	yellow = "\033[33m"
	red    = "\033[31m"
	cyan   = "\033[36m"
	gray   = "\033[90m"
)

// Helper untuk warna status
func statusColor(code int) string {
	switch {
	case code >= 200 && code < 300:
		return green
	case code >= 300 && code < 400:
		return cyan
	case code >= 400 && code < 500:
		return yellow
	default:
		return red
	}
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// 1. GUNAKAN WRAPPER BAWAAN CHI
		// Ini menggantikan struct `responseWriter` manual yang kamu buat sebelumnya.
		// Wrapper ini otomatis menangkap Status Code dan Bytes Written.
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		// 2. Jalankan handler selanjutnya
		// Penting: Masukkan `ww` (wrapper), bukan `w` asli.
		next.ServeHTTP(ww, r)

		// 3. Ambil Status Code dari wrapper Chi
		status := ww.Status()

		// 4. Print Output (Sesuai format kode kamu)
		fmt.Printf("%s%s%s %s %s%d%s %s%.2fms%s\n",
			cyan, r.Method, reset,
			r.URL.Path,
			statusColor(status), status, reset,
			gray, float64(time.Since(start).Microseconds())/1000, reset,
		)
	})
}
