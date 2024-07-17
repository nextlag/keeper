// Package gzip - middleware gzip
package gzip

import (
	"net/http"
	"strings"
)

// MwGzip returns middleware to handle gzip compression for HTTP requests and responses.
func MwGzip() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ow := w

			acceptEncoding := r.Header.Get("Accept-Encoding")
			supportGzip := strings.Contains(acceptEncoding, "gzip")

			if supportGzip {
				cw := NewCompressWriter(w)
				ow = cw
				defer cw.Close()
			}

			contentEncoding := r.Header.Get("Content-Encoding")
			sendGzip := strings.Contains(contentEncoding, "gzip")

			if sendGzip {
				cr, err := NewCompressReader(r.Body)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				r.Body = cr
				defer cr.Close()
			}

			next.ServeHTTP(ow, r)
		}
		return http.HandlerFunc(fn)
	}
}
