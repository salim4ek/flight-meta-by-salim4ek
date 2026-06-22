// Package httpapi exposes the search service over HTTP/JSON.
package httpapi

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"flightmeta/internal/search"
	"flightmeta/internal/sources"
)

// Server wires the HTTP handlers to the search orchestrator.
type Server struct {
	orch        *search.Orchestrator
	log         *slog.Logger
	allowOrigin string
}

// New returns an http.Handler with all routes and middleware configured.
func New(orch *search.Orchestrator, log *slog.Logger, allowOrigin string) http.Handler {
	s := &Server{orch: orch, log: log, allowOrigin: allowOrigin}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", s.handleHealth)
	mux.HandleFunc("GET /search", s.handleSearch)
	return s.middleware(mux)
}

// middleware adds minimal security headers and CORS for the React dev origin.
// Full hardening (CSP, rate limiting) lands in Phase 7.
func (s *Server) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("X-Frame-Options", "DENY")
		h.Set("Referrer-Policy", "no-referrer")
		if s.allowOrigin != "" {
			h.Set("Access-Control-Allow-Origin", s.allowOrigin)
			h.Set("Vary", "Origin")
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) handleSearch(w http.ResponseWriter, r *http.Request) {
	q, err := parseQuery(r)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, s.orch.Search(r.Context(), q))
}

// parseQuery validates and normalizes the search query string. All input is
// validated here so downstream code can trust it.
func parseQuery(r *http.Request) (sources.Query, error) {
	v := r.URL.Query()

	origin := strings.ToUpper(strings.TrimSpace(v.Get("origin")))
	dest := strings.ToUpper(strings.TrimSpace(v.Get("destination")))
	if !isIATA(origin) || !isIATA(dest) {
		return sources.Query{}, fmt.Errorf("origin and destination must be 3-letter IATA codes")
	}
	if origin == dest {
		return sources.Query{}, fmt.Errorf("origin and destination must differ")
	}

	depart, err := time.Parse("2006-01-02", v.Get("depart"))
	if err != nil {
		return sources.Query{}, fmt.Errorf("depart must be YYYY-MM-DD")
	}
	var ret *time.Time
	if rs := strings.TrimSpace(v.Get("return")); rs != "" {
		rt, err := time.Parse("2006-01-02", rs)
		if err != nil {
			return sources.Query{}, fmt.Errorf("return must be YYYY-MM-DD")
		}
		if rt.Before(depart) {
			return sources.Query{}, fmt.Errorf("return must not be before depart")
		}
		ret = &rt
	}

	pax := 1
	if ps := strings.TrimSpace(v.Get("passengers")); ps != "" {
		n, err := strconv.Atoi(ps)
		if err != nil || n < 1 || n > 9 {
			return sources.Query{}, fmt.Errorf("passengers must be 1..9")
		}
		pax = n
	}

	stops := sources.StopsMode(strings.TrimSpace(v.Get("stops")))
	switch stops {
	case sources.StopsAny, sources.StopsDirect, sources.StopsOne, sources.StopsOnePlus:
	default:
		return sources.Query{}, fmt.Errorf("stops must be one of: direct, one, one_plus")
	}

	passport := strings.ToUpper(strings.TrimSpace(v.Get("passport")))
	if passport == "" {
		passport = "RU"
	}

	return sources.Query{
		Origin:      origin,
		Destination: dest,
		DepartDate:  depart,
		ReturnDate:  ret,
		Passengers:  pax,
		Passport:    passport,
		Currency:    strings.ToUpper(strings.TrimSpace(v.Get("currency"))),
		Filters: sources.Filters{
			Stops:               stops,
			IncludeAirlines:     splitCSVUpper(v.Get("airlines")),
			ExcludeAirlines:     splitCSVUpper(v.Get("exclude_airlines")),
			AllowSelfTransfer:   v.Get("self_transfer") != "false", // allowed by default
			OnlyVisaFreeTransit: v.Get("visa_free_transit") == "true",
		},
	}, nil
}

func isIATA(s string) bool {
	if len(s) != 3 {
		return false
	}
	for _, c := range s {
		if c < 'A' || c > 'Z' {
			return false
		}
	}
	return true
}

func splitCSVUpper(s string) []string {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p = strings.ToUpper(strings.TrimSpace(p)); p != "" {
			out = append(out, p)
		}
	}
	return out
}

func writeJSON(w http.ResponseWriter, code int, body any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(body)
}
