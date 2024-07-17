package gzip

import (
	"compress/gzip"
	"io"
	"net/http"
)

// CompressWriter provides a gzip-compressed HTTP response writer.
type CompressWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

// NewCompressWriter creates a new CompressWriter that wraps an existing http.ResponseWriter.
func NewCompressWriter(w http.ResponseWriter) *CompressWriter {
	return &CompressWriter{
		w:  w,
		zw: gzip.NewWriter(w),
	}
}

// Header returns the header map that will be sent by the CompressWriter.
func (c *CompressWriter) Header() http.Header {
	return c.w.Header()
}

// Write writes compressed data to the underlying gzip.Writer.
func (c *CompressWriter) Write(p []byte) (int, error) {
	return c.zw.Write(p)
}

// WriteHeader sends an HTTP response header with the provided status code.
// If the status code is less than 300, it sets the "Content-Encoding" header to "gzip".
func (c *CompressWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		c.w.Header().Set("Content-Encoding", "gzip")
	}
	c.w.WriteHeader(statusCode)
}

// Close closes the underlying gzip.Writer.
func (c *CompressWriter) Close() error {
	return c.zw.Close()
}

// CompressReader provides a gzip-decompressed reader for HTTP request bodies.
type CompressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

// NewCompressReader creates a new CompressReader that wraps an existing io.ReadCloser.
// It assumes the input data is gzip-compressed and initializes a gzip.Reader.
func NewCompressReader(r io.ReadCloser) (*CompressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	return &CompressReader{
		r:  r,
		zr: zr,
	}, nil
}

// Read reads decompressed data from the underlying gzip.Reader.
func (c CompressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

// Close closes the underlying gzip.Reader and the original io.ReadCloser.
func (c *CompressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}
