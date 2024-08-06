package request

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"

	"github.com/nextlag/keeper/pkg/logger/l"
)

// Request contains the request fields for the logger.
type Request struct {
	Method      string `json:"method"`
	Path        string `json:"path"`
	RemoteAddr  string `json:"remote_addr"`
	UserAgent   string `json:"user_agent"`
	RequestID   string `json:"request_id"`
	ContentType string `json:"content_type,omitempty"`
	Status      int    `json:"status"`
	Bytes       int    `json:"bytes,omitempty"`
	Duration    string `json:"duration"`
	Compress    string `json:"compress,omitempty"`
}

// MwRequest creates middleware for logging HTTP requests.
func MwRequest(log *l.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			requestFields := Request{
				Method:      r.Method,
				Path:        r.URL.Path,
				RemoteAddr:  r.RemoteAddr,
				UserAgent:   r.UserAgent(),
				RequestID:   middleware.GetReqID(r.Context()),
				ContentType: r.Header.Get("Content-Type"),
			}

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			t1 := time.Now()

			defer func() {
				requestFields.Status = ww.Status()
				requestFields.Bytes = ww.BytesWritten()
				requestFields.Duration = time.Since(t1).String()
				requestFields.Compress = ww.Header().Get("Content-Encoding")

				if requestFields.Status >= http.StatusInternalServerError {
					log.Error("request completed with error", "error logger", requestFields)
				} else {
					log.Debug("request", "request fields", requestFields)
				}
			}()

			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}
