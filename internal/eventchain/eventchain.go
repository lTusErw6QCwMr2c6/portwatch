// Package eventchain links related events into causal chains using a
// parent/child relationship tracked by event ID.
package eventchain

import (
	"sync"
	"time"

	"github.com/user/portwatch/internal/alert"
)

// Link represents a parent→child relationship between two events.
type Link struct {
	ParentID  string
	ChildID   string
	LinkedAt  time.Time
}

// Chain tracks causal links between events.
type Chain struct {
	mu      sync.RWMutex
	parents map[string]string   // childID → parentID
	children map[string][]string // parentID → []childID
}

// New returns an empty Chain.
func New() *Chain {
	return &Chain{
		parents:  make(map[string]string),
		children: make(map[string][]string),
	}
}

// Link records that child was caused by parent.
func (c *Chain) Link(parent, child alert.Event) {
	pID := parent.Port.String()
	cID := child.Port.String()
	c.mu.Lock()
	defer c.mu.Unlock()
	c.parents[cID] = pID
	c.children[pID] = append(c.children[pID], cID)
}

// Parent returns the parent event key for a given child key, if any.
func (c *Chain) Parent(child alert.Event) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	p, ok := c.parents[child.Port.String()]
	return p, ok
}

// Children returns all child event keys for a given parent event.
func (c *Chain) Children(parent alert.Event) []string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	kids := c.children[parent.Port.String()]
	out := make([]string, len(kids))
	copy(out, kids)
	return out
}

// Reset clears all recorded links.
func (c *Chain) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.parents = make(map[string]string)
	c.children = make(map[string][]string)
}
