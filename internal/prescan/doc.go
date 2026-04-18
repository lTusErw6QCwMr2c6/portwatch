// Package prescan provides a startup warm-up scan that captures the initial
// state of open ports on the host. The result is used to seed the watcher's
// baseline, preventing spurious open/close events on the first diff cycle.
package prescan
