package httputils

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type gzipRW struct {
	http.ResponseWriter
	sniffed bool
	gz      io.Writer
}

func (w *gzipRW) Write(b []byte) (int, error) {
	if !w.sniffed {
		if w.Header().Get("Content-Type") == "" {
			w.Header().Set("Content-Type", http.DetectContentType(b))
		}
		w.sniffed = true
	}
	return w.gz.Write(b)
}

func GzipHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			fn(w, r)
			return
		}
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)

		defer gz.Close()
		gzrw := gzipRW{
			ResponseWriter: w,
			gz:             gz,
		}
		fn(&gzrw, r)
	}
}
