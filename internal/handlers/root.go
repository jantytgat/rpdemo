package handlers

import (
	"encoding/base64"
	"log/slog"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/jantytgat/go-kit/pkg/slogd"
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
		Method:          r.Method,
		Uri:             r.RequestURI,
		BackgroundColor: h.bgColor,
		TextColor:       "white",
	}

	fullHeader := r.Header.Get("X-Ns-Fullheader")
	var fullHeaderBytes []byte
	var originalHeaders []string
	var originalHeadersMap = make(map[string]string)
	for _, k := range keys {
		originalHeadersMap[k] = ""
	}
	if fullHeaderBytes, err = base64.StdEncoding.DecodeString(fullHeader); err != nil {
		fullHeader = "Failed to decode"
	} else {
		fullHeader = string(fullHeaderBytes)
		originalHeaders = strings.Split(fullHeader, "\n")
		for _, header := range originalHeaders {
			if cKey, cValue, found := strings.Cut(header, ":"); found {
				originalHeadersMap[cKey] = cValue
			}
		}
	}

	headerData := HeaderData{
		Headers:         headers,
		OriginalHeaders: originalHeadersMap,
		Keys:            keys,
		Timestamp:       time.Now(),
	}

	main(info, headerData).Render(r.Context(), w)
	slogd.Logger().LogAttrs(r.Context(), slogd.LevelInfo, "root handler", slog.Time("timestamp", headerData.Timestamp))
}
