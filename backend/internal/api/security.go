package api

import (
	"net/http"
	"strings"
)

// contentSecurityPolicy limits what the browser may load for the
// application.
//
// The only foreign script is the SheetJS build that reads an Excel file, so
// its CDN is the single extra source. The styles of the single file
// components are injected by Vue at runtime, which needs 'unsafe-inline' for
// the style rule only.
const contentSecurityPolicy = "default-src 'self'; " +
	"base-uri 'self'; " +
	"object-src 'none'; " +
	"frame-ancestors 'none'; " +
	"form-action 'self'; " +
	"script-src 'self' https://cdn.sheetjs.com; " +
	"style-src 'self' 'unsafe-inline'; " +
	"img-src 'self' data: blob:; " +
	"font-src 'self' data:; " +
	"connect-src 'self' blob:"

// SecurityHeaders adds the response headers that protect the application.
//
// The headers are set on every response, including the static files of the
// frontend, so they are also applied when the application is served by the
// Go server itself instead of a reverse proxy.
func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("X-Frame-Options", "DENY")
		h.Set("Referrer-Policy", "no-referrer")
		h.Set("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
		h.Set("Content-Security-Policy", contentSecurityPolicy)
		if strings.HasPrefix(r.URL.Path, "/api/") {
			// The answers of the API carry the data of the people who
			// registered. They must not stay in a cache of the browser or of
			// a proxy that sits in front of the application.
			h.Set("Cache-Control", "no-store")
		}
		if IsSecureRequest(r) {
			// Only a secure response may carry the instruction, otherwise a
			// browser would pin the plain host to HTTPS.
			h.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}
		next.ServeHTTP(w, r)
	})
}

// IsSecureRequest reports whether a request arrived over TLS.
//
// A request that arrives over plain HTTP but carries the forwarding header
// of a proxy that terminated TLS counts as secure as well.
func IsSecureRequest(r *http.Request) bool {
	if r == nil {
		return false
	}
	if r.TLS != nil {
		return true
	}
	proto := r.Header.Get("X-Forwarded-Proto")
	if i := strings.IndexByte(proto, ','); i >= 0 {
		proto = proto[:i]
	}
	return strings.EqualFold(strings.TrimSpace(proto), "https")
}

// withCORS allows the given origins to call the API.
//
// An origin that is not on the list gets no permission at all, so a page of
// another site cannot use the session of the user. An empty list means only
// the same origin may call the API, which is the normal case in production
// because the frontend and the API are served by the same host.
func withCORS(next http.Handler, allowed []string) http.Handler {
	allowed = normalizeOrigins(allowed)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" && originAllowed(origin, allowed) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
			w.Header().Set("Vary", "Origin")
		}
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// normalizeOrigins trims the origins and drops the empty ones.
func normalizeOrigins(origins []string) []string {
	out := make([]string, 0, len(origins))
	for _, o := range origins {
		o = strings.TrimSpace(o)
		if o == "" {
			continue
		}
		out = append(out, strings.TrimSuffix(o, "/"))
	}
	return out
}

// originAllowed reports whether an origin may call the API.
func originAllowed(origin string, allowed []string) bool {
	origin = strings.TrimSuffix(strings.TrimSpace(origin), "/")
	for _, a := range allowed {
		if strings.EqualFold(a, origin) {
			return true
		}
	}
	return false
}
