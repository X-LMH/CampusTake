package logxext

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"
)

const (
	colorReset  = "\033[0m"
	colorGreen  = "\033[32m"
	colorCyan   = "\033[36m"
	colorYellow = "\033[33m"
	colorRed    = "\033[31m"

	// 背景色
	bgGreen   = "\033[42m"
	bgYellow  = "\033[43m"
	bgRed     = "\033[41m"
	bgCyan    = "\033[46m"
	bgBlue    = "\033[44m"
	bgMagenta = "\033[45m"
	bgWhite   = "\033[47m"
	bgBlack   = "\033[40m"

	// 前景色
	fgBlack = "\033[30m"
	fgWhite = "\033[97m"
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
	// 为常见 HTTP 方法选择背景色与合适的前景色，保证可读性
	switch method {
	case "GET":
		// 绿色背景，黑字
		return bgGreen + fgBlack + method + colorReset
	case "POST":
		// 蓝色背景，白字（常用于创建）
		return bgBlue + fgWhite + method + colorReset
	case "PUT":
		// 青色背景，黑字（常用于更新/替换）
		return bgCyan + fgBlack + method + colorReset
	case "DELETE":
		// 红色背景，白字（危险/删除）
		return bgRed + fgWhite + method + colorReset
	case "PATCH":
		// 品红背景，白字（部分更新）
		return bgMagenta + fgWhite + method + colorReset
	case "OPTIONS":
		// 黄色背景，黑字（信息性）
		return bgYellow + fgBlack + method + colorReset
	case "HEAD":
		// 白色背景，黑字（无正文）
		return bgWhite + fgBlack + method + colorReset
	case "TRACE", "CONNECT":
		// 黑色背景，白字（不常用 / 特殊）
		return bgBlack + fgWhite + method + colorReset
	default:
		// 默认青色背景黑字
		return bgCyan + fgBlack + method + colorReset
	}
}

func HttpLoggerMiddleware() func(next http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			rw := &statusRecorder{ResponseWriter: w}
			next(rw, r)

			cost := time.Since(start)
			ms := float64(cost.Nanoseconds()) / 1e6
			// 固定宽度显示毫秒，保证列对齐（例如 "  2.32ms" / " 16.15ms"）
			costStr := fmt.Sprintf("%6.2fms", ms)

			// 优先从代理头获取真实 IP
			ip := r.Header.Get("X-Real-Ip")
			if ip == "" {
				if xf := r.Header.Get("X-Forwarded-For"); xf != "" {
					parts := strings.Split(xf, ",")
					ip = strings.TrimSpace(parts[0])
				}
			}
			if ip == "" {
				ip = r.RemoteAddr
				if host, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
					ip = host
				}
			}

			ts := time.Now().Format("2006-01-02 15:04:05")

			statusCode := rw.status
			if statusCode == 0 {
				statusCode = http.StatusOK
			}

			method := strings.ToUpper(r.Method)

			// 将方法字段填充到固定宽度，保证后续路径列对齐
			const methodWidth = 7 // 最大方法长度（如 "CONNECT"）
			pad := methodWidth - len(method)
			if pad < 0 {
				pad = 0
			}
			methodField := colorMethod(method) + strings.Repeat(" ", pad)

			// 移除行尾字节数显示，简化输出
			line := fmt.Sprintf(
				"%s | %s | %s | %s | %s %s",
				ts,
				colorStatus(statusCode),
				costStr,
				ip,
				methodField,
				r.URL.Path,
			)

			fmt.Println(line)
		}
	}
}

type statusRecorder struct {
	http.ResponseWriter
	status int
	bytes  int64
}

func (r *statusRecorder) Header() http.Header {
	return r.ResponseWriter.Header()
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func (r *statusRecorder) Write(b []byte) (int, error) {
	if r.status == 0 {
		// 如果未显式调用 WriteHeader，则认为是 200
		r.status = http.StatusOK
	}
	n, err := r.ResponseWriter.Write(b)
	r.bytes += int64(n)
	return n, err
}

// 转发可选接口，保证兼容性
func (r *statusRecorder) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hj, ok := r.ResponseWriter.(http.Hijacker); ok {
		return hj.Hijack()
	}
	return nil, nil, fmt.Errorf("Hijacker not supported")
}

func (r *statusRecorder) Flush() {
	if f, ok := r.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

func (r *statusRecorder) CloseNotify() <-chan bool {
	if cn, ok := r.ResponseWriter.(http.CloseNotifier); ok {
		return cn.CloseNotify()
	}
	// 返回一个永不关闭的通道，避免 nil 使用
	return make(chan bool)
}

func (r *statusRecorder) Push(target string, opts *http.PushOptions) error {
	if p, ok := r.ResponseWriter.(http.Pusher); ok {
		return p.Push(target, opts)
	}
	return http.ErrNotSupported
}
