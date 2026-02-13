package main

import (
	"fmt"
	"net"
	"net/http"
	"strings"
)

func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For first
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// Take the first IP in the comma-separated list
		ip := strings.TrimSpace(strings.Split(xff, ",")[0])
		if ip != "" {
			return ip
		}
	}

	// Check X-Real-IP
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return strings.TrimSpace(xri)
	}

	// Fall back to RemoteAddr, stripping port
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		// RemoteAddr might not have a port
		return r.RemoteAddr
	}
	return host
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "https://ryou.onl")

	if r.URL.Path == "/health" {
		if r.Method == http.MethodGet {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "ok")
			return
		}
		http.NotFound(w, r)
		return
	}

	if r.URL.Path == "/" {
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.WriteHeader(http.StatusNoContent)
			return
		}
		if r.Method == http.MethodGet {
			ip := getClientIP(r)
			w.Header().Set("Content-Type", "text/plain")
			fmt.Fprint(w, ip)
			return
		}
		http.NotFound(w, r)
		return
	}

	http.NotFound(w, r)
}

func main() {
	fmt.Println("ip-service listening on :8080")
	http.HandleFunc("/", handler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("server error: %v\n", err)
	}
}
