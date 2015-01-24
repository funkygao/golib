// The breaker package provides circuit breaker primitives to enable one to
// safely interact with potentially-unreliable subsystems in a way that allows
// a graceful alternative behavior while in degraded state.
//
// For instance, if your server depends on an external database and it becomes
// unavailable, your server can use a breaker to track these failures and
// provide an alternative workflow during this period of unavailability.
//
// A circuit has two states:
//
// CLOSED: The system decorated with this breaker is assumed to be available
// and the dependents thereof may use it freely.
//
// OPEN: The system decorated with this breaker is assumed to be unavailable
// and the dependents thereof should not use it at this time.
package breaker

import (
	"fmt"
	"sync"
	"time"
)

type state int

const (
	closed state = iota
	open
)

func (s state) String() string {
	switch s {
	case open:
		return "OPEN"

	case closed:
		return "CLOSED"

	default:
		return "InvalidState"
	}

}

type instantProvider interface {
	Now() time.Time
}

type systemInstantProvider struct{}

func (i *systemInstantProvider) Now() time.Time {
	return time.Now()
}

var (
	systemTimer = &systemInstantProvider{}
)

// Consecutive provides a simple circuit breaker with a deadline expiration
// strategy to determine when to reclose itself.  Its initial state is CLOSED
// and will only open once a threshold of consecutive failures is reached.
//
// It is concurrency safe.
type Consecutive struct {
	// RetryTimeout is added to the time when the circuit first opens to determine
	// the next time the decorating operation will be allowed to occur.
	RetryTimeout time.Duration

	// FailureAllowance is the maximum consecutive failures that may occur before
	// the circuit will open.
	FailureAllowance uint

	mutex sync.RWMutex

	state    state
	failures uint

	nextClose       time.Time
	instantProvider instantProvider
}

func (b *Consecutive) String() string {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	return fmt.Sprintf("Consecutive Breaker %s with %d Allowance and %s Deadline and %d Failures", b.state, b.FailureAllowance, b.RetryTimeout, b.failures)
}

func (b *Consecutive) enabled() bool {
	return b.FailureAllowance > 0
}

// Fail marks the decorated subsystem as having an operation fail and may
// trigger its subsequent circuit opening.
//
// This should only be called if you can divine that the underlying failure
// is or should be considered transient for the given domain of work.  For
// instance, high latency, erroneous responses, unavailability count as
// so-called temporary failures.  Things like invalid user queries count as
// permanent errors and would never make sense to be retried.
func (b *Consecutive) Fail() {
	if !b.enabled() {
		return
	}

	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.failures++

	if b.state == open {
		return
	}

	if b.instantProvider == nil {
		b.instantProvider = systemTimer
	}

	if b.failures > b.FailureAllowance {
		b.nextClose = b.instantProvider.Now().Add(b.RetryTimeout)
		b.state = open
	}
}

// Succeed marks the decorated subsystem as having an operation succeed and will
// its circuits closure if it's presently open.
func (b *Consecutive) Succeed() {
	if !b.enabled() {
		return
	}

	b.mutex.Lock()
	defer b.mutex.Unlock()

	if b.state == open {
		b.reset()
	}
}

// Open indicates whether the circuit for this subsystem is presently open.
func (b *Consecutive) Open() bool {
	if !b.enabled() {
		return false
	}

	b.mutex.Lock()
	defer b.mutex.Unlock()

	if b.instantProvider == nil {
		b.instantProvider = systemTimer
	}

	switch {
	case b.state == closed:
		return false
	case b.nextClose.Before(b.instantProvider.Now()):
		b.reset()
		return false
	default:
		return true
	}

}

// Reset returns this circuit back to its default state: closed.
//
// This should be called if you can divine that the underlying subsystem has
// become unavailable before the deadline threshold has been reached.
func (b *Consecutive) Reset() {
	if !b.enabled() {
		return
	}

	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.reset()
}

func (b *Consecutive) reset() {
	b.failures = 0
	b.state = closed
}
