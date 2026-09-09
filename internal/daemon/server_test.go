//go:build linux

package daemon

import (
	"testing"
	"time"
)

func TestShouldIdleShutdown_PendingNotificationBlocksShutdown(t *testing.T) {
	s := &Server{
		focusCtx:     map[uint32]focusInfo{1: {}},
		lastNotifID:  map[string]uint32{},
		idleTimeout:  5 * time.Minute,
		lastActivity: time.Now().Add(-time.Hour), // long past the timeout
	}

	if s.shouldIdleShutdown() {
		t.Fatal("expected shutdown to be blocked while a notification is still pending a click")
	}

	// Pending activity should have reset the idle clock too.
	s.activityMu.Lock()
	idle := time.Since(s.lastActivity)
	s.activityMu.Unlock()
	if idle >= s.idleTimeout {
		t.Fatalf("expected lastActivity to be refreshed, got idle=%v", idle)
	}
}

func TestShouldIdleShutdown_NoPendingRespectsTimeout(t *testing.T) {
	s := &Server{
		focusCtx:     map[uint32]focusInfo{},
		lastNotifID:  map[string]uint32{},
		idleTimeout:  5 * time.Minute,
		lastActivity: time.Now().Add(-time.Hour),
	}

	if !s.shouldIdleShutdown() {
		t.Fatal("expected shutdown once idle timeout elapsed with no pending notifications")
	}
}

func TestShouldIdleShutdown_NoPendingWithinTimeout(t *testing.T) {
	s := &Server{
		focusCtx:     map[uint32]focusInfo{},
		lastNotifID:  map[string]uint32{},
		idleTimeout:  5 * time.Minute,
		lastActivity: time.Now(),
	}

	if s.shouldIdleShutdown() {
		t.Fatal("expected no shutdown while still within idle timeout")
	}
}
