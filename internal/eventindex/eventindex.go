package eventindex

import (
	"sync"
	"time"

	"github.com/user/portwatch/internal/alert"
)

// Index provides fast lookup of events by port, protocol, and type.
type Index struct {
	mu      sync.RWMutex
	byPort  map[int][]entry
	byProto map[string][]entry
	byType  map[string][]entry
	ttl     time.Duration
}

type entry struct {
	event   alert.Event
	addedAt time.Time
}

// New returns an Index with the given TTL for eviction.
func New(ttl time.Duration) *Index {
	if ttl <= 0 {
		ttl = 5 * time.Minute
	}
	return &Index{
		byPort:  make(map[int][]entry),
		byProto: make(map[string][]entry),
		byType:  make(map[string][]entry),
		ttl:     ttl,
	}
}

// Add indexes the given event.
func (idx *Index) Add(ev alert.Event) {
	idx.mu.Lock()
	defer idx.mu.Unlock()
	e := entry{event: ev, addedAt: time.Now()}
	port := ev.Port.Number
	proto := ev.Port.Protocol
	kind := string(ev.Kind)
	idx.byPort[port] = append(idx.byPort[port], e)
	idx.byProto[proto] = append(idx.byProto[proto], e)
	idx.byType[kind] = append(idx.byType[kind], e)
}

// ByPort returns all live events matching the given port number.
func (idx *Index) ByPort(port int) []alert.Event {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	return idx.collect(idx.byPort[port])
}

// ByProtocol returns all live events matching the given protocol.
func (idx *Index) ByProtocol(proto string) []alert.Event {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	return idx.collect(idx.byProto[proto])
}

// ByType returns all live events matching the given event kind.
func (idx *Index) ByType(kind string) []alert.Event {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	return idx.collect(idx.byType[kind])
}

// Evict removes all entries older than the configured TTL.
func (idx *Index) Evict() {
	idx.mu.Lock()
	defer idx.mu.Unlock()
	now := time.Now()
	idx.byPort = evictMap(idx.byPort, now, idx.ttl)
	idx.byProto = evictMap(idx.byProto, now, idx.ttl)
	idx.byType = evictMap(idx.byType, now, idx.ttl)
}

func (idx *Index) collect(entries []entry) []alert.Event {
	now := time.Now()
	out := make([]alert.Event, 0, len(entries))
	for _, e := range entries {
		if now.Sub(e.addedAt) <= idx.ttl {
			out = append(out, e.event)
		}
	}
	return out
}

func evictMap(m map[string][]entry, now time.Time, ttl time.Duration) map[string][]entry {
	result := make(map[string][]entry, len(m))
	for k, entries := range m {
		var live []entry
		for _, e := range entries {
			if now.Sub(e.addedAt) <= ttl {
				live = append(live, e)
			}
		}
		if len(live) > 0 {
			result[k] = live
		}
	}
	return result
}
