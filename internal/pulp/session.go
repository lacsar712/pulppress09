package pulp

import "sync"

type LiveSession struct {
	mu    sync.Mutex
	label string
	n     int
}

func NewLiveSession(label string) *LiveSession {
	return &LiveSession{label: label}
}

func (s *LiveSession) Inc() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.n++
}

func (s *LiveSession) Value() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.n
}

// CloneSession forks an independent LiveSession that carries the same nameplate
// (label) AND the same pressure-walk sequence number (n) as the source at the
// moment of the split. Both fields are copied under the source's lock so the
// snapshot is consistent.
//
// Why the full state must travel with the clone:
// n is the cross-request, monotone running state of the nip pressure walk
// (every Inc() advances it by one). label alone is just the static nameplate.
// Before this fix CloneSession only copied label and left n=0, so a forked
// (e.g. day-shift) copy restarted the walk from scratch while the original
// (night-shift) session kept advancing the real counter. After the fork the two
// sessions held divergent counters for the same physical nip: requests routed
// to the clone under-pressured from a baseline of zero and fell out of step
// with the DCS beat, while the original's n kept moving on. Copying n makes the
// clone resume from the parent's last-observed sequence number, so the only
// subsequent drift is the intended, independent incrementing each session does
// on its own requests — no more phantom reset-to-zero at handover.
func CloneSession(s *LiveSession) *LiveSession {
	s.mu.Lock()
	defer s.mu.Unlock()
	return &LiveSession{label: s.label, n: s.n}
}
