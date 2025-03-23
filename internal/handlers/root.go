package handlers

import (
	"net/http"
	"sort"
	"time"
)

func NewRootHandler() RootHandler {
	return RootHandler{}
}

type RootHandler struct{}

func (h RootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	headers := make(map[string]string)

	var keys []string
	for k := range r.Header {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		headers[k] = r.Header.Get(k)
	}

	rootHandler(keys, headers, time.Now()).Render(r.Context(), w)
	// w.WriteHeader(http.StatusOK)
	//
	// var keys []string
	// for k := range r.Header {
	// 	keys = append(keys, k)
	// }
	// sort.Strings(keys)
	//
	// for _, key := range keys {
	// 	w.Write([]byte(key))
	// 	w.Write([]byte("\t"))
	// 	for _, value := range r.Header.Values(key) {
	// 		w.Write([]byte(value))
	// 		w.Write([]byte("\r\n"))
	// 	}
	// 	w.Write([]byte("\r\n"))
	// }
}
