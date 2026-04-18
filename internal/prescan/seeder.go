package prescan

import (
	"fmt"

	"github.com/user/portwatch/internal/snapshot"
)

// Seeder seeds a snapshot store from a pre-scan result so the first watcher
// cycle has a valid baseline and does not emit noise events.
type Seeder struct {
	store *snapshot.Store
}

// NewSeeder creates a Seeder backed by the provided store.
func NewSeeder(store *snapshot.Store) *Seeder {
	return &Seeder{store: store}
}

// Seed saves a new snapshot derived from the pre-scan Result.
func (s *Seeder) Seed(r Result) error {
	if s.store == nil {
		return fmt.Errorf("prescan: seeder has nil store")
	}
	snap := snapshot.New(r.Ports)
	return s.store.Save(snap)
}
