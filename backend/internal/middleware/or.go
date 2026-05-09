package middleware

import (
	"bytes"
	"net/http"
)

type bufferedResponse struct {
	header http.Header
	code   int
	body   bytes.Buffer
}

func (r *bufferedResponse) Header() http.Header { return r.header }

func (r *bufferedResponse) WriteHeader(code int) { r.code = code }

func (r *bufferedResponse) Write(b []byte) (int, error) {
	if r.code == 0 {
		r.code = http.StatusOK
	}
	return r.body.Write(b)
}

// Or returns middleware that runs first against a buffered response writer.
// If first responds with 401 the buffer is discarded and second runs against
// the real writer. Any other response from first is flushed to the real writer.
func Or(first, second func(http.Handler) http.Handler) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rec := &bufferedResponse{header: make(http.Header)}
			first(next).ServeHTTP(rec, r)
			if rec.code == http.StatusUnauthorized {
				second(next).ServeHTTP(w, r)
				return
			}
			for k, v := range rec.header {
				w.Header()[k] = v
			}
			if rec.code != 0 {
				w.WriteHeader(rec.code)
			}
			w.Write(rec.body.Bytes()) //nolint:errcheck
		})
	}
}
