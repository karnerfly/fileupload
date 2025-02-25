package middlewares

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type GzipResponseWriter struct {
	http.ResponseWriter
	io.Writer
}

func GzipMiddleware(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Set("Content-Encoding", "gzip")

		gw := gzip.NewWriter(w)
		defer gw.Close()

		next.ServeHTTP(GzipResponseWriter{ResponseWriter: w, Writer: gw}, r)
	})
}

func (grw GzipResponseWriter) Write(p []byte) (int, error) {
	return grw.Writer.Write(p)
}
