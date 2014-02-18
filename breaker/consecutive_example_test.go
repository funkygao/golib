package breaker_test

import (
	breaker "."
	"fmt"
	"time"
)

// RemoteService is a fictitious interface to an encapsulation of some
// unreliable subsystem that is outside of our domain of control.  It should
// be treated for this example as merely a recording of behavior.
type RemoteService struct {
	// The subsystem's circuit breaker.
	Breaker breaker.Consecutive
	// The pre-recorded success or failure results that are used to drive
	// the behavior of this example.
	Results []bool
}

// ConductRequest is the supposed interface point for user's of this subsystem.
// They call it as necessary to perform whatever work to yield the result they
// want.
func (s *RemoteService) ConductRequest() {
	// For purposes of not convoluting the example, we use real sleep operations
	// here.
	time.Sleep(time.Second/2 + time.Second/4)

	// If the circuit is broken, merely bail.
	if s.Breaker.Open() {
		fmt.Println("Unavailable; Trying Again Later...")

		return
	}

	// Emulate the actual remote interface here that is supposedly unreliable.
	err := s.performRequest()
	// WARNING: We make an implicit assumption that any err value is retryable
	// and not a permanent error.
	if err != nil {
		fmt.Println("Operation Failed")
		s.Breaker.Fail()
	} else {
		fmt.Println("Operation Succeeded")
		s.Breaker.Succeed()
	}
}

// performRequest models an interaction with an unreliable external
// system---e.g., a remote API server.
func (s *RemoteService) performRequest() error {
	result := s.Results[0]
	s.Results = s.Results[1:]

	if !result {
		return fmt.Errorf("Temporary Unavailable")
	}

	return nil
}

func ExampleConsecutive() {
	subsystem := &RemoteService{
		Breaker: breaker.Consecutive{
			FailureAllowance: 2,
			RetryTimeout:     time.Second,
		},
		Results: []bool{
			true,  // Success.
			true,  // Success.
			false, // One-off failure; do not trip circuit.
			true,  // This success negates past failures.
			false, // String of contiguous failures to create open circuit.
			false, // Open.  :-(
			false, // Open; however, we've timed out.
			true,  // We have a success here.
		},
	}

	subsystem.ConductRequest()
	subsystem.ConductRequest()
	subsystem.ConductRequest()
	subsystem.ConductRequest()
	subsystem.ConductRequest()
	subsystem.ConductRequest()
	subsystem.ConductRequest()
	subsystem.ConductRequest()
	subsystem.ConductRequest()
	// Output:
	// Operation Succeeded
	// Operation Succeeded
	// Operation Failed
	// Operation Succeeded
	// Operation Failed
	// Operation Failed
	// Unavailable; Trying Again Later...
	// Operation Failed
	// Operation Succeeded
}
