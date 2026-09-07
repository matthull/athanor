package tmux

import "time"

// WakePane triggers SIGWINCH by doing a resize dance on the target window.
// When Claude runs in a detached tmux session, its TUI library may not process
// stdin until a terminal event occurs. This simulates attaching by resizing.
//
// Uses resize-window (not resize-pane) because resize-pane doesn't work on
// single-pane sessions — the pane already fills the window. [G:tmux.go:1226-1257]
func WakePane(r *Runner, target string) {
	width, err := r.GetWindowWidth(target)
	if err != nil || width < 1 {
		return // session may be dead
	}
	// Bump width +1, wait 50ms, restore
	_ = r.ResizeWindow(target, width+1)
	time.Sleep(50 * time.Millisecond)
	_ = r.ResizeWindow(target, width)

	// Reset window-size to "latest" — resize-window sets it to "manual"
	// which permanently locks dimensions and causes dot borders when
	// a human later attaches with a different terminal size.
	// [G:tmux.go:1251-1256]
	_ = r.SetWindowOption(target, "window-size", "latest")
}

// WakePaneIfDetached triggers SIGWINCH only if the session is detached.
// Skip the wake on attached sessions — Claude is already processing terminal
// events, and the resize dance adds unnecessary latency. [G:tmux.go:1262-1267]
func WakePaneIfDetached(r *Runner, target string) {
	if r.IsSessionAttached(target) {
		return
	}
	WakePane(r, target)
}
