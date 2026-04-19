// Package eventid generates unique, sortable identifiers for port events.
package eventid

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync/atomic"
	"time"
)

var counter uint64

// ID is a unique event identifier.
type ID struct {
	Timestamp time.Time
	Sequence  uint64
	Random    string
}

// New returns a new unique ID.
func New() ID {
	seq := atomic.AddUint64(&counter, 1)
	b := make([]byte, 4)
	_, _ = rand.Read(b)
	return ID{
		Timestamp: time.Now().UTC(),
		Sequence:  seq,
		Random:    hex.EncodeToString(b),
	}
}

// String returns a string representation of the ID.
func (id ID) String() string {
	return fmt.Sprintf("%d-%06d-%s", id.Timestamp.UnixMilli(), id.Sequence, id.Random)
}

// Before reports whether id occurred before other.
func (id ID) Before(other ID) bool {
	if id.Timestamp.Equal(other.Timestamp) {
		return id.Sequence < other.Sequence
	}
	return id.Timestamp.Before(other.Timestamp)
}

// Reset resets the internal sequence counter. Intended for testing only.
func Reset() {
	atomic.StoreUint64(&counter, 0)
}
