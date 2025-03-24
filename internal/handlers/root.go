package handlers

import (
	"encoding/base64"
	"net/http"
	"sort"
	"time"
)

func NewRootHandler(colorFlag string) RootHandler {
	return RootHandler{
		bgColor: colorFlag,
	}
}

type RootHandler struct {
	bgColor string
}

func (h RootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error
	headers := make(map[string]string)

	var keys []string
	for k := range r.Header {
		if k != "X-Ns-Fullheader" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	for _, k := range keys {
		headers[k] = r.Header.Get(k)
	}

	info := PageInfo{
		Title:           "Reverse Proxy Demo",
		Language:        "en",
		BackgroundColor: h.bgColor,
		TextColor:       "white",
	}

	fullHeader := r.Header.Get("X-Ns-Fullheader")
	var fullHeaderBytes []byte
	if fullHeaderBytes, err = base64.StdEncoding.DecodeString(fullHeader); err != nil {
		fullHeader = "Failed to decode"
	} else {
		fullHeader = string(fullHeaderBytes)
	}

	headerData := HeaderData{
		Headers:         headers,
		OriginalHeaders: fullHeader,
		Keys:            keys,
		Timestamp:       time.Now(),
	}

	main(info, headerData).Render(r.Context(), w)
}
