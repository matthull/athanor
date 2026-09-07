package tmux

import (
	"sync"
	"time"
)

// LockTimeout is the default duration to wait for a per-target send lock.
// Prevents a hung tmux from permanently blocking all future sends to a target.
// [G:tmux.go:35]
const LockTimeout = 30 * time.Second

// targetLocks holds per-target channel semaphores. Uses sync.Map for
// lock-free concurrent access. [G:tmux.go:29]
var targetLocks sync.Map // map[string]chan struct{}

// getTargetSem returns the channel semaphore for a target, creating one
// if it doesn't exist. The semaphore is a buffered channel of size 1 —
// send to acquire, receive to release. [G:tmux.go:1183-1187]
func getTargetSem(target string) chan struct{} {
	sem := make(chan struct{}, 1)
	actual, _ := targetLocks.LoadOrStore(target, sem)
	return actual.(chan struct{})
}

// AcquireLock attempts to acquire the per-target send lock within timeout.
// Returns true if acquired, false if the timeout expired. Channel-based
// (not mutex) so callers can give up instead of blocking forever.
// [G:tmux.go:1189-1199]
func AcquireLock(target string, timeout time.Duration) bool {
	sem := getTargetSem(target)
	select {
	case sem <- struct{}{}:
		return true
	case <-time.After(timeout):
		return false
	}
}

// ReleaseLock releases the per-target send lock. The default case prevents
// blocking if called without holding the lock. [G:tmux.go:1201-1209]
func ReleaseLock(target string) {
	sem := getTargetSem(target)
	select {
	case <-sem:
	default:
		// Lock wasn't held — shouldn't happen, but don't block
	}
}
