package httplogger

import (
	"log"
	"net"
	"net/http"
	"time"
)

type stResponseWriter struct {
	http.ResponseWriter
	HTTPStatus   int
	ResponseSize int
}

func (w *stResponseWriter) WriteHeader(status int) {
	w.HTTPStatus = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *stResponseWriter) Flush() {
	z := w.ResponseWriter
	if f, ok := z.(http.Flusher); ok {
		f.Flush()
	}
}

func (w *stResponseWriter) CloseNotify() <-chan bool {
	z := w.ResponseWriter
	return z.(http.CloseNotifier).CloseNotify()
}

func (w *stResponseWriter) Write(b []byte) (int, error) {
	if w.HTTPStatus == 0 {
		w.HTTPStatus = 200
	}
	w.ResponseSize = len(b)
	return w.ResponseWriter.Write(b)
}

func HTTPLogger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		interceptWriter := stResponseWriter{w, 0, 0}

		host, _, _ := net.SplitHostPort(r.RemoteAddr)

		handler.ServeHTTP(&interceptWriter, r)
		log.Printf("HTTP - %s - - %s \"%s %s %s\" %d %d \"%s\" \"%s\" %dus\n",
			host,
			t.Format("[02/Jan/2006:15:04:05 -0700]"),
			r.Method,
			r.URL.Path,
			r.Proto,
			interceptWriter.HTTPStatus,
			interceptWriter.ResponseSize,
			r.Referer(),
			r.UserAgent(),
			time.Since(t)/1000,
		)
	})
}
