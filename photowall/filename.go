package photowall

import (
	"strconv"
	"sync"
)

// NamerFunc can be used instead of the Namer type
type NamerFunc func(p Photo) string

// Name transforms NamerFunc to a Namer
func (n NamerFunc) Name(p Photo) string {
	return n(p)
}

// Namer generates unique names
type Namer interface {
	Name(Photo) string
}

// DateNamer represents name generator based on date/time
type DateNamer struct {
	last    string
	counter int
	format  string
	mutex   sync.Mutex
}

// NewDateNamer generates a new name generator based on date/time
func NewDateNamer(format string) *DateNamer {
	return &DateNamer{
		format: format,
	}
}

// Name generates an unique date-based string. Format defines date format e.g. "2006-01-02_150405".
// Conflicting names are resolved by adding "_1..".
// Thread-safe
func (n *DateNamer) Name(p Photo) string {
	n.mutex.Lock()
	nowStr := p.CreatedAt().Format(n.format)
	if nowStr == n.last {
		n.counter++
		nowStr = nowStr + "_" + strconv.Itoa(n.counter)
	} else {
		n.last = nowStr
		n.counter = 0
	}
	n.mutex.Unlock()
	return nowStr
}
