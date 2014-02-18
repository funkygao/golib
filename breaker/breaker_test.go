package breaker

import (
	"testing"
	"time"
)

type testInstantProvider []time.Time

func (i *testInstantProvider) Now() time.Time {
	defer func() {
		*i = (*i)[1:]
	}()

	return (*i)[0]
}

func TestConsecutiveLifecycle(t *testing.T) {
	consecutive := Consecutive{}
	if consecutive.Open() {
		t.Fatal("expected false, got true")
	}

	consecutive.Reset()
	if consecutive.Open() {
		t.Fatal("expected false, got true")
	}

}

func TestConsecutive(t *testing.T) {
	now := time.Now()

	type in struct {
		successes   []bool
		consecutive Consecutive
	}

	type out struct {
		open []bool
	}

	var scenarios = []struct {
		in  in
		out out
	}{
		// All successes.
		{
			in: in{
				consecutive: Consecutive{
					RetryTimeout:     time.Minute,
					FailureAllowance: 5,
				},
				successes: []bool{
					true,
					true,
					true,
					true,
					true,
				},
			},
			out: out{
				open: []bool{
					false,
					false,
					false,
					false,
					false,
				},
			},
		},
		// All failures under allowance.
		{
			in: in{
				consecutive: Consecutive{
					RetryTimeout:     time.Minute,
					FailureAllowance: 5,
				},
				successes: []bool{
					false,
					false,
					false,
					false,
					false,
				},
			},
			out: out{
				open: []bool{
					false,
					false,
					false,
					false,
					false,
				},
			},
		},
		// All failures over allowance.
		{
			in: in{
				consecutive: Consecutive{
					RetryTimeout:     time.Minute,
					FailureAllowance: 5,
				},
				successes: []bool{
					false,
					false,
					false,
					false,
					false,
					false,
				},
			},
			out: out{
				open: []bool{
					false,
					false,
					false,
					false,
					false,
					true,
				},
			},
		},
		// All failures over allowance and continues.
		{
			in: in{
				consecutive: Consecutive{
					RetryTimeout:     time.Minute,
					FailureAllowance: 5,
					// instantProvider: &testInstantProvider{
					// 	now,
					// 	now.Add(80 * time.Second),
					// 	// now.Add(40 * time.Second),
					// 	// now.Add(60 * time.Second),
					// 	// now.Add(80 * time.Second),
					// },
				},
				successes: []bool{
					false,
					false,
					false,
					false,
					false,
					false,
					false,
					false,
					false,
					false,
				},
			},
			out: out{
				open: []bool{
					false,
					false,
					false,
					false,
					false,
					true,
					true,
					true,
					true,
					true,
				},
			},
		},
		// All failures over allowance but expires.
		{
			in: in{
				consecutive: Consecutive{
					RetryTimeout:     time.Minute,
					FailureAllowance: 5,
					instantProvider: &testInstantProvider{
						now,
						now,
						now.Add(80 * time.Second),
					},
				},
				successes: []bool{
					false,
					false,
					false,
					false,
					false,
					false,
					false,
					false,
					false,
					false,
				},
			},
			out: out{
				open: []bool{
					false,
					false,
					false,
					false,
					false,
					true,
					false,
					false,
					false,
					false,
				},
			},
		},
	}

	for i, scenario := range scenarios {
		in := scenario.in
		out := scenario.out
		consecutive := in.consecutive

		if len(in.successes) != len(out.open) {
			t.Fatalf("%d. expected %d inputs, got %d", i, len(out.open), len(in.successes))
		}

		for j, success := range in.successes {
			if success {
				consecutive.Succeed()
			} else {
				consecutive.Fail()
			}

			open := consecutive.Open()

			if open != out.open[j] {
				t.Fatalf("%d.%d. expected state %v, got %v", i, j, out.open[j], open)
			}
		}
	}
}
