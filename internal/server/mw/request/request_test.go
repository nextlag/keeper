package request_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	config "github.com/nextlag/keeper/config/server"
	"github.com/nextlag/keeper/internal/server/mw/request"
	"github.com/nextlag/keeper/pkg/logger/l"
)

func TestLoggerMiddleware(t *testing.T) {
	tests := []struct {
		name     string
		method   string
		path     string
		expected string
	}{
		{
			name:     "GET request",
			method:   http.MethodGet,
			path:     "/test",
			expected: "GET /test",
		},
		{
			name:     "POST request",
			method:   http.MethodPost,
			path:     "/data",
			expected: "POST /data",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cfg, err := config.Load()
			if err != nil {
				require.Error(t, err)
			}
			testLogger := l.NewLogger(cfg)
			mw := request.MwRequest(testLogger)

			handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tc.method, r.Method)
				assert.Equal(t, tc.path, r.URL.Path)

				w.WriteHeader(http.StatusOK)
				_, err = w.Write([]byte("OK"))
				if err != nil {
					return
				}
			}))

			req := httptest.NewRequest(tc.method, tc.path, nil)
			req.Header.Set("Content-Type", "application/json")
			req.RemoteAddr = "127.0.0.1:12345"
			req.Header.Set("User-Agent", "Test User Agent")
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)
			assert.Equal(t, http.StatusOK, rec.Code)
		})
	}
}

func BenchmarkNewLoggerMiddleware(b *testing.B) {
	cfg, err := config.Load()
	if err != nil {
		require.Error(b, err)
	}
	testLogger := l.NewLogger(cfg)
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mw := request.MwRequest(testLogger)
		rr := &fakeResponseWriter{}
		mw(handler).ServeHTTP(rr, req)
	}
}

type fakeResponseWriter struct {
	header http.Header
	code   int
}

func (rw *fakeResponseWriter) Header() http.Header {
	if rw.header == nil {
		rw.header = make(http.Header)
	}
	return rw.header
}

func (rw *fakeResponseWriter) Write(b []byte) (int, error) {
	return len(b), nil
}

func (rw *fakeResponseWriter) WriteHeader(statusCode int) {
	rw.code = statusCode
}
