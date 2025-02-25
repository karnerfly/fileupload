package middlewares

import (
	"compress/gzip"
	"net/http"
	"strings"
)

type GzipResponseWriter struct {
	rw         http.ResponseWriter
	gzipWriter *gzip.Writer
}

func NewGzipResponseWriter(w http.ResponseWriter) *GzipResponseWriter {
	gw := gzip.NewWriter(w)
	return &GzipResponseWriter{
		rw:         w,
		gzipWriter: gw,
	}
}

func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		grw := NewGzipResponseWriter(w)
		defer grw.Flush()

		grw.Header().Add("Content-Encoding", "gzip")

		next.ServeHTTP(grw, r)
	})
}

func (grw *GzipResponseWriter) Write(b []byte) (int, error) {
	return grw.gzipWriter.Write(b)
}

func (grw *GzipResponseWriter) Header() http.Header {
	return grw.rw.Header()
}

func (grw *GzipResponseWriter) WriteHeader(statusCode int) {
	grw.rw.WriteHeader(statusCode)
}

func (grw *GzipResponseWriter) Flush() {
	grw.gzipWriter.Flush()
	grw.gzipWriter.Close()
}
