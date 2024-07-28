package main

import (
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

var AddrRegex = regexp.MustCompile(`^(([a-zA-Z]|[a-zA-Z][a-zA-Z0-9\-]*[a-zA-Z0-9])\.)*([A-Za-z]|[A-Za-z][A-Za-z0-9\-]*[A-Za-z0-9])$`)

func traceRouteHandler(w http.ResponseWriter, r *http.Request) {
	path, _, _ := strings.Cut(r.URL.Path, "?")
	trimmedPath := strings.Trim(path, "/")
	addr := trimmedPath[strings.LastIndex(trimmedPath, "/")+1:]

	if !AddrRegex.MatchString(addr) {
		WriteBadRequest(w, "Invalid addr")
		return
	}

	maxHops, err := strconv.Atoi(defaultString(r.URL.Query().Get("maxHops"), "30"))
	if err != nil {
		WriteBadRequest(w, "maxHops must be an integer")
		return
	}

	hops := TraceRoute(addr, maxHops)

	result := make([]*TraceHopResponse, len(hops))
	for i, hop := range hops {
		result[i] = hop.toTraceHopResponse(i)
	}

	response, err := WriteJSON(w, result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = saveToRedis(GenerateRedisKey(addr, maxHops), response)
	if err != nil {
		log.Fatal(err)
	}
}

func RunTraceRouteServer(listen string) {
	handler := &RegexpHandler{}
	handler.HandleFunc(regexp.MustCompile(`^/trace/[^/]+$`), traceRouteHandler)

	log.Printf("Listening on %s\n", listen)
	err := http.ListenAndServe(listen, handler)
	if err != nil {
		log.Fatalf("Failed to start server. %v", err)
	}
}
