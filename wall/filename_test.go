package wall

import (
	"testing"
	"time"
)

func TestDateNamerUnique(t *testing.T) {
	const format = "2006-01-02_150405"
	var n Namer
	p := NewPhoto("", 0, 0, "", time.Now())
	n = NewDateNamer(format)
	t1 := n.Name(p)
	t2 := n.Name(p)
	if len(t1) < len(format) {
		t.Fatalf("Length of name to small: %s (%d)", t1, len(t1))
	}
	if t1 == t2 {
		t.Fatalf("Namer generates non-unique names: %s", t1)
	}
}
