package tmux

import (
	"fmt"
	"strings"
	"time"
)

// Timing constants for the 8-step send protocol. These values are empirically
// derived from production — do not change without evidence.
const (
	sendKeysChunkSize  = 512                    // [G:tmux.go:1314] bytes per chunk
	postSendDelay      = 500 * time.Millisecond // [G:tmux.go:1453] empirically required
	escEnterGap        = 600 * time.Millisecond // [G:tmux.go:1460] must exceed readline keyseq-timeout (500ms)
	enterRetryDelay    = 200 * time.Millisecond // [G:tmux.go:1471] between Enter retries
	enterMaxAttempts   = 3                      // [G:tmux.go:1469]
	copyModeCancelWait = 50 * time.Millisecond  // [G:tmux.go:1440]
	interChunkDelay    = 10 * time.Millisecond  // [G:tmux.go:1340]
	nudgeRetryInterval = 500 * time.Millisecond // [G:constants:54] base retry backoff
	maxRetryInterval   = 2 * time.Second        // [G:tmux.go:1386] backoff cap
	defaultSendTimeout = 15 * time.Second       // default timeout for send-keys retry
)

// SendOpts configures the behavior of Send.
type SendOpts struct {
	SkipEscape bool          // omit Escape + 600ms delay
	Timeout    time.Duration // max time for send-keys retry (default 15s)
	LockHeld   bool          // skip lock acquisition (caller already holds the lock)
}

// Send implements the 8-step send protocol for reliable message delivery
// to a tmux target. This is the core of whisper — every send executes
// these steps. [G:tmux.go:1421-1482]
func (r *Runner) Send(target, message string, opts SendOpts) error {
	timeout := opts.Timeout
	if timeout == 0 {
		timeout = defaultSendTimeout
	}

	// Step 1: Acquire send lock [G:tmux.go:1424]
	if !opts.LockHeld {
		if !AcquireLock(target, LockTimeout) {
			return fmt.Errorf("send lock timeout for %q: previous send may be hung", target)
		}
		defer ReleaseLock(target)
	}

	// Step 2: Exit copy mode [G:tmux.go:1438]
	inMode, err := r.InCopyMode(target)
	if err == nil && inMode {
		_ = r.CancelCopyMode(target)
		time.Sleep(copyModeCancelWait)
	}

	// Step 3: Sanitize [G:tmux.go:1444]
	sanitized := Sanitize(message)

	// Step 4: Send text with chunking + retry [G:tmux.go:1448]
	if err := r.sendMessageToTarget(target, sanitized, timeout); err != nil {
		return err
	}

	// Step 5: Post-send delay [G:tmux.go:1453]
	time.Sleep(postSendDelay)

	// Step 6: Send Escape [G:tmux.go:1458]
	if !opts.SkipEscape {
		_ = r.SendKeys(target, "Escape")

		// Step 7: ESC→Enter gap [G:tmux.go:1460]
		// Bash readline's keyseq-timeout defaults to 500ms. ESC + Enter within
		// that window becomes M-Enter (meta-return), which does nothing.
		// 600ms ensures ESC is processed alone before Enter arrives.
		time.Sleep(escEnterGap)
	}

	// Step 8: Send Enter with retry [G:tmux.go:1469-1481]
	var lastErr error
	for attempt := 0; attempt < enterMaxAttempts; attempt++ {
		if attempt > 0 {
			time.Sleep(enterRetryDelay)
		}
		if err := r.SendKeys(target, "Enter"); err != nil {
			lastErr = err
			continue
		}
		// Step 8b: Wake detached sessions
		WakePaneIfDetached(r, target)
		return nil
	}
	return fmt.Errorf("failed to send Enter after %d attempts: %w", enterMaxAttempts, lastErr)
}

// sendMessageToTarget sends text to a tmux target, chunking messages larger
// than sendKeysChunkSize bytes. The first chunk uses retry logic to handle
// the startup race; subsequent chunks use plain SendKeysLiteral since if the
// first chunk worked, the target is ready. [G:tmux.go:1316-1344]
func (r *Runner) sendMessageToTarget(target, text string, timeout time.Duration) error {
	if len(text) <= sendKeysChunkSize {
		return r.sendKeysLiteralWithRetry(target, text, timeout)
	}
	for i := 0; i < len(text); i += sendKeysChunkSize {
		end := i + sendKeysChunkSize
		if end > len(text) {
			end = len(text)
		}
		chunk := text[i:end]
		if i == 0 {
			if err := r.sendKeysLiteralWithRetry(target, chunk, timeout); err != nil {
				return err
			}
		} else {
			if err := r.SendKeysLiteral(target, chunk); err != nil {
				return err
			}
		}
		if end < len(text) {
			time.Sleep(interChunkDelay)
		}
	}
	return nil
}

// sendKeysLiteralWithRetry sends literal text to a tmux target, retrying on
// transient errors with exponential backoff: 500ms → 750ms → 1125ms → 1687ms → 2s (cap).
// Non-transient errors (session gone, no server) fail immediately.
// [G:tmux.go:1359-1391]
func (r *Runner) sendKeysLiteralWithRetry(target, text string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	interval := nudgeRetryInterval
	var lastErr error

	for time.Now().Before(deadline) {
		err := r.SendKeysLiteral(target, text)
		if err == nil {
			return nil
		}
		if !isTransientSendKeysError(err) {
			return err // non-transient — fail fast
		}
		lastErr = err
		remaining := time.Until(deadline)
		if remaining <= 0 {
			break
		}
		sleep := interval
		if sleep > remaining {
			sleep = remaining
		}
		time.Sleep(sleep)
		// Exponential backoff: 500ms → 750ms → 1125ms → 1687ms → 2s (cap)
		interval = interval * 3 / 2
		if interval > maxRetryInterval {
			interval = maxRetryInterval
		}
	}
	return fmt.Errorf("target not ready for input after %s: %w", timeout, lastErr)
}

// isTransientSendKeysError returns true if the error is transient and safe to
// retry. Only "not in a mode" is retryable — this happens when the TUI hasn't
// initialized its input handling during cold startup. [G:tmux.go:1272-1278]
func isTransientSendKeysError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "not in a mode")
}
