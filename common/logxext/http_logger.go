package logxext

import (
	"fmt"
	"net"
	"net/http"
	"time"
)

const (
	colorReset  = "\033[0m"
	colorGreen  = "\033[32m"
	colorCyan   = "\033[36m"
	colorYellow = "\033[33m"
	colorRed    = "\033[31m"

	// 背景色
	bgGreen  = "\033[42m"
	bgYellow = "\033[43m"
	bgRed    = "\033[41m"
	bgCyan   = "\033[46m"

	// 前景色
	fgBlack = "\033[30m"
)

func colorStatus(code int) string {
	switch {
	case code >= 500:
		return bgRed + fgBlack + fmt.Sprint(code) + colorReset
	case code >= 400:
		return bgYellow + fgBlack + fmt.Sprint(code) + colorReset
	default:
		return bgGreen + fgBlack + fmt.Sprint(code) + colorReset
	}
}

func colorMethod(method string) string {
	return bgCyan + fgBlack + method + colorReset
}

func HttpLoggerMiddleware() func(next http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			rw := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
			next(rw, r)

			cost := time.Since(start)

			ip := r.RemoteAddr
			if host, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
				ip = host
			}

			ts := time.Now().Format("2006-01-02 15:04:05")

			line := fmt.Sprintf(
				"%s | %s | %s | %s | %s %s",
				ts,
				colorStatus(rw.status),
				cost.String(),
				ip,
				colorMethod(r.Method),
				r.URL.Path,
			)

			fmt.Println(line)
		}
	}
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}
