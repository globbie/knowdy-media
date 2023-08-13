package webgate

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"
)

func (wg *WebGate) authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		strArr := strings.Split(authHeader, " ")
		if len(strArr) != 2 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
                token := strArr[1]

                // TODO
		userId := token
		/*userId, err := wg.AccountUseCases.Authenticate(token)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}*/
		ctx := context.WithValue(r.Context(), userIdContextKey, userId)
		next(w, r.WithContext(ctx))
	}
}

type responseWriterObserver struct {
	http.ResponseWriter
	status int
	wroteHeader bool
}

func (o *responseWriterObserver) WriteHeader(code int) {
	o.ResponseWriter.WriteHeader(code)
	if o.wroteHeader {
		return
	}
	o.wroteHeader = true
	o.status = code
}

func (o *responseWriterObserver) StatusCode() int {
	if !o.wroteHeader {
		return http.StatusOK
	}
	return o.status
}

func (wg *WebGate) logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		o := &responseWriterObserver{ResponseWriter: w}
		next.ServeHTTP(o, r)
		fmt.Printf("method: %s; url: %s; status-code: %d; remote-addr: %s; duration: %v;\n",
			r.Method, r.URL.String(), o.StatusCode(), r.RemoteAddr, time.Since(start))
	})
}
