package request

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nextlag/keeper/config/server"
	"github.com/nextlag/keeper/pkg/logger/l"
)

func TestMwRequest(t *testing.T) {
	cfg, err := server.Load()
	require.NoError(t, err)

	log := l.NewLogger(cfg)

	req := httptest.NewRequest(http.MethodGet, "/test/request", nil)
	req = req.WithContext(context.WithValue(req.Context(), middleware.RequestIDKey, "test-request-id"))
	rr := httptest.NewRecorder()

	requestMiddleware := MwRequest(log)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte("test response"))
		if err != nil {
			return
		}
	})

	wrappedHandler := requestMiddleware(handler)
	wrappedHandler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "test response", rr.Body.String())

}
