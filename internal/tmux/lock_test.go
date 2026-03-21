package tmux

import (
	"testing"
	"time"
)

func TestAcquireRelease(t *testing.T) {
	t.Parallel()
	target := "test-lock-basic"
	if !AcquireLock(target, time.Second) {
		t.Fatal("should acquire uncontested lock")
	}
	ReleaseLock(target)
	// Should be able to re-acquire after release
	if !AcquireLock(target, time.Second) {
		t.Fatal("should re-acquire after release")
	}
	ReleaseLock(target)
}

func TestAcquireLockTimeout(t *testing.T) {
	t.Parallel()
	target := "test-lock-timeout"
	if !AcquireLock(target, time.Second) {
		t.Fatal("first acquire should succeed")
	}
	defer ReleaseLock(target)

	// Second acquire should timeout
	start := time.Now()
	if AcquireLock(target, 50*time.Millisecond) {
		t.Fatal("should not acquire contested lock")
	}
	elapsed := time.Since(start)
	if elapsed < 40*time.Millisecond {
		t.Errorf("timeout too fast: %v", elapsed)
	}
}

func TestDifferentTargetsIndependent(t *testing.T) {
	t.Parallel()
	target1 := "test-lock-ind-1"
	target2 := "test-lock-ind-2"
	if !AcquireLock(target1, time.Second) {
		t.Fatal("should acquire target1")
	}
	defer ReleaseLock(target1)
	// target2 should be independently acquirable
	if !AcquireLock(target2, time.Second) {
		t.Fatal("should acquire target2 independently")
	}
	ReleaseLock(target2)
}

func TestReleaseWithoutAcquire(t *testing.T) {
	t.Parallel()
	// Should not panic or block
	ReleaseLock("test-lock-never-held")
}
