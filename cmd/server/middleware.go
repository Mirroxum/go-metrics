package main

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
)

func DecompressGzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Encoding") == "gzip" {
			gz, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, "Failed to decompress gzip body", http.StatusBadRequest)
				return
			}
			defer gz.Close()

			var buf bytes.Buffer
			if _, err := io.Copy(&buf, gz); err != nil {
				http.Error(w, "Failed to read decompressed body", http.StatusInternalServerError)
				return
			}

			r.Body = io.NopCloser(&buf)
			r.ContentLength = int64(buf.Len())
		}

		next.ServeHTTP(w, r)
	})
}
