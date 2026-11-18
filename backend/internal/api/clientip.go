package api

import (
	"net"
	"net/http"
	"strings"
)

// clientAddress returns the address a request comes from.
//
// The address is used to count the attempts of a login and of a
// registration, so it must not be in the hands of the caller. A header such
// as X-Forwarded-For can be sent by anyone, so it is only read when the
// application is told that its own reverse proxy writes the real address
// there (TrustProxy). Otherwise the address of the connection is used, which
// no caller can change.
//
// Behind a reverse proxy without TrustProxy every request shares the address
// of the proxy, so the counters of the addresses become the counters of the
// whole application. Set TRUST_PROXY=true in that setup.
func (s *Server) clientAddress(r *http.Request) string {
	if r == nil {
		return ""
	}
	if s.cfg.TrustProxy {
		if addr := lastForwardedFor(r); addr != "" {
			return addr
		}
	}
	return remoteAddress(r)
}

// lastForwardedFor returns the last address of the X-Forwarded-For header.
//
// The header is a list of the addresses a request passed through. The proxy
// of this application appends the address it saw to that list, so the last
// entry is the one the proxy observed and not one the caller could choose.
func lastForwardedFor(r *http.Request) string {
	header := r.Header.Get("X-Forwarded-For")
	if header == "" {
		return ""
	}
	parts := strings.Split(header, ",")
	for i := len(parts) - 1; i >= 0; i-- {
		if addr := strings.TrimSpace(parts[i]); addr != "" {
			return addr
		}
	}
	return ""
}

// remoteAddress returns the host part of the address of the connection.
func remoteAddress(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return strings.TrimSpace(r.RemoteAddr)
	}
	return host
}
