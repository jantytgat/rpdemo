package handlers

import (
	"net/http"
	"sort"
)

func NewRootHandler() RootHandler {
	return RootHandler{}
}

type RootHandler struct{}

func (h RootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	var keys []string
	for k := range r.Header {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		w.Write([]byte(key))
		w.Write([]byte("\t"))
		for _, value := range r.Header.Values(key) {
			w.Write([]byte(value))
			w.Write([]byte("\r\n"))
		}
		w.Write([]byte("\r\n"))
	}
}
