package model

import (
	"crypto/sha256"
	"fmt"
	"strings"
)

// LegacyUserID builds a stable id for an account of the previous
// application, which does not carry one.
//
// The id has the shape of a UUID, so it can be told apart from the id of an
// account that lives only in this application.
func LegacyUserID(username string) string {
	sum := sha256.Sum256([]byte("user:" + strings.ToLower(strings.TrimSpace(username))))
	b := sum[:16]
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}
