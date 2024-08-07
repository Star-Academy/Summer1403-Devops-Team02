package main

import (
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

var DomainRegex = regexp.MustCompile(`^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\-]*[A-Za-z0-9])$`)
var AddrRegex = regexp.MustCompile(`^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`)

func traceRouteHandler(w http.ResponseWriter, r *http.Request) {
	path, _, _ := strings.Cut(r.URL.Path, "?")
	trimmedPath := strings.Trim(path, "/")
	addr := trimmedPath[strings.LastIndex(trimmedPath, "/")+1:]

	if !AddrRegex.MatchString(addr) && !DomainRegex.MatchString(addr) {
		WriteBadRequest(w, "Invalid addr "+addr)
		return
	}

	maxHops, err := strconv.Atoi(defaultString(r.URL.Query().Get("maxHops"), "30"))
	if err != nil {
		WriteBadRequest(w, "maxHops must be an integer")
		return
	}

	hops, err := TraceRoute(addr, maxHops)
	if err != nil {
		WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}

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
		log.Println(err)
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
