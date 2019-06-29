/*!
 * go-rs/compression
 * Copyright(c) 2019 Roshan Gade
 * MIT Licensed
 */
package compression

import (
	"compress/gzip"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"

	"github.com/go-rs/rest-api-framework"
)

var pool sync.Pool

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

type compressionHandler struct {
	api *rest.API
}

func (h *compressionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		h.api.ServeHTTP(w, r)
		return
	}

	w.Header().Set("Content-Encoding", "gzip")

	gz := pool.Get().(*gzip.Writer)
	defer pool.Put(gz)

	gz.Reset(w)
	defer func() {
		_ = gz.Flush()
		_ = gz.Close()
	}()

	h.api.ServeHTTP(&gzipResponseWriter{ResponseWriter: w, Writer: gz}, r)
}

func Handler(api *rest.API, level int) http.Handler {
	if level == 0 {
		level = gzip.DefaultCompression
	}

	pool = sync.Pool{
		New: func() interface{} {
			w, err := gzip.NewWriterLevel(ioutil.Discard, level)
			if err != nil {
				panic(err)
			}
			return w
		},
	}
	return &compressionHandler{api}
}
