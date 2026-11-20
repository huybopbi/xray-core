package shadowsocks

import (
	"testing"
	"time"
)

func TestAttackDefenseCloseStopsCleanupLoop(t *testing.T) {
	d := NewAttackDefense(&DefenseConfig{
		MaxFailures:     20,
		BanDuration:     time.Minute,
		CleanupInterval: time.Hour,
	})

	errc := make(chan error, 1)
	go func() {
		errc <- d.Close()
	}()

	select {
	case err := <-errc:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Close did not return; cleanupLoop is still running")
	}

	if err := d.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestServerCloseStopsAttackDefense(t *testing.T) {
	v := new(Validator)
	v.initializeIfNeeded()
	s := &Server{validator: v}
	if v.attackDefense == nil {
		t.Fatal("expected attack defense to start with the validator")
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
}
