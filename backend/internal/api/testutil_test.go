package api

import (
	"net/http/cookiejar"
)

// cookieJar returns a jar that keeps the session cookie across requests.
func cookieJar() (*cookiejar.Jar, error) {
	return cookiejar.New(nil)
}
