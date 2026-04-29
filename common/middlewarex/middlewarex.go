package middlewarex

import (
	"bytes"
	"context"
	"io"
	"net/http"
)

type ctxKey string

const RequestBodyKey ctxKey = "request_body"

func CacheRequestBody(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Body != nil {
			bodyBytes, _ := io.ReadAll(r.Body)
			_ = r.Body.Close()

			// 重新放回去，避免后续绑定参数时读不到
			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

			ctx := context.WithValue(r.Context(), RequestBodyKey, string(bodyBytes))
			r = r.WithContext(ctx)
		}

		next.ServeHTTP(w, r)
	})
}

func GetRequestBody(r *http.Request) string {
	v := r.Context().Value(RequestBodyKey)
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}
