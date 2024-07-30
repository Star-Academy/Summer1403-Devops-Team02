package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"
)

func defaultString(s, def string) string {
	if s == "" {
		return def
	}
	return s
}

type route struct {
	pattern *regexp.Regexp
	handler http.Handler
}

type RegexpHandler struct {
	routes []*route
}

func (h *RegexpHandler) Handler(pattern *regexp.Regexp, handler http.Handler) {
	h.routes = append(h.routes, &route{pattern, handler})
}

func (h *RegexpHandler) HandleFunc(pattern *regexp.Regexp, handler func(http.ResponseWriter, *http.Request)) {
	h.routes = append(h.routes, &route{pattern, http.HandlerFunc(handler)})
}

func (h *RegexpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, route := range h.routes {
		if route.pattern.MatchString(r.URL.Path) {
			route.handler.ServeHTTP(w, r)
			return
		}
	}
	// no pattern matched; send 404 response
	http.NotFound(w, r)
}

func WriteBadRequest(w http.ResponseWriter, msg string) []byte {
	resp := WriteError(w, msg, http.StatusBadRequest)
	log.Printf("Bad request: %s\n", resp)
	return resp
}

func WriteError(w http.ResponseWriter, msg string, code int) []byte {
	w.WriteHeader(code)
	resp, _ := json.Marshal(map[string]string{"error": msg})

	log.Printf("%s: %s\n", http.StatusText(code), resp)

	w.Write(resp)
	return resp
}

func WriteJSON(w http.ResponseWriter, data interface{}) ([]byte, error) {
	w.Header().Set("Content-Type", "application/json")
	jsonData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, err
	}

	w.Write(jsonData)

	return jsonData, nil
}

func GenerateRedisKey(addr string, maxHops int) string {
	timestamp := time.Now().Unix()
	return fmt.Sprintf("%s:%d:%d", addr, maxHops, timestamp)
}
