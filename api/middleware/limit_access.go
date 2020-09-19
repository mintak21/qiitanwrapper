package middleware

import (
	"net/http"

	"github.com/didip/tollbooth"
)

const (
	maxBodyByteSize = 100 * 1024 * 1024 // 100MB
	requestQuota    = 10                // 10 request per second
)

// BodySizeLimit restricts access of body size
func BodySizeLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, maxBodyByteSize)
		next.ServeHTTP(w, r)
	})
}

// RequestQuotaLimit restricts 1 request/second
func RequestQuotaLimit(next http.Handler) http.Handler {
	limiter := tollbooth.NewLimiter(requestQuota, nil).
		SetIPLookups([]string{"RemoteAddr", "X-Forwarded-For", "X-Real-IP"}).
		SetMethods([]string{"GET", "POST"}).
		SetMessage("reached maximum request quota")
	return tollbooth.LimitHandler(limiter, next)
}
