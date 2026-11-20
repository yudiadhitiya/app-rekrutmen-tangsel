package supabase

import (
	"strings"
	"testing"

	"rekrutmen-tangsel/backend/internal/model"
)

// A value that is put into a filter as it is may end that filter and start
// another one, so a caller could ask for something else than the one asked
// for. The characters that change the meaning of a filter are removed.
func TestPgSafeRemovesTheGrammarOfAFilter(t *testing.T) {
	for _, dangerous := range []string{
		"a,b",
		"a)or(b",
		"a.b",
		`a"b`,
		"a*b",
		"a%b",
		"a:b",
	} {
		got := pgSafe(dangerous)
		if strings.ContainsAny(got, `,()."*%:`) {
			t.Errorf("pgSafe(%q) = %q, masih memuat karakter filter", dangerous, got)
		}
	}
}

// The words that a name and a number may hold survive the cleaning.
func TestPgSafeKeepsTheUsualCharacters(t *testing.T) {
	got := pgSafe("Budi Santoso-01")
	if got != "Budi Santoso-01" {
		t.Errorf("pgSafe = %q, ingin %q", got, "Budi Santoso-01")
	}
}

// The filter that reaches PostgREST holds no character that could end it.
func TestListRekrutmenSanitizesTheFilter(t *testing.T) {
	f := newFakePostgREST()
	ts := f.serve(t)
	defer ts.Close()

	st := New(ts.URL, "test-key", "")
	if _, _, err := st.ListRekrutmen(model.Filter{Q: "a,or(created_by.eq.x)"}); err != nil {
		t.Fatalf("ListRekrutmen: %v", err)
	}

	filter := f.lastQuery.Get("or")
	if filter == "" {
		t.Fatal("filter pencarian tidak terkirim")
	}
	// The injected part must not reach the server as a filter of its own.
	if strings.Contains(filter, "created_by.eq") {
		t.Errorf("filter = %q, masih memuat filter yang disuntikkan", filter)
	}
	// Only the parentheses of the group itself may be there.
	if strings.Count(filter, "(") != 1 || strings.Count(filter, ")") != 1 {
		t.Errorf("filter = %q, tanda kurung tidak sesuai", filter)
	}
}
