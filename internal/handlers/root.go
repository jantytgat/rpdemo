package handlers

import (
	"net/http"
)

func NewRootHandler() RootHandler {
	return RootHandler{}
}

type RootHandler struct{}

func (h RootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello World\r\n"))
	for key, header := range r.Header {
		for _, value := range header {
			w.Write([]byte(key))
			w.Write([]byte("\t"))
			w.Write([]byte(value))
			w.Write([]byte("\r\n"))
		}
	}
}
