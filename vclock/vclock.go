package vclock

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

type VectorClock struct {
	Values []int
}

func (c VectorClock) Equals(other VectorClock) bool {
	if len(c.Values) != len(other.Values) {
		return false
	}
	for i, v := range c.Values {
		if v != other.Values[i] {
			return false
		}
	}
	return true
}

func (c VectorClock) happensBefore(other VectorClock) bool {
	oneLessThan := false
	allLessOrEqual := true
	for i, v := range c.Values {
		allLessOrEqual = allLessOrEqual && v <= other.Values[i]
		oneLessThan = oneLessThan || v < other.Values[i]
	}
	return oneLessThan && allLessOrEqual
}

func (c VectorClock) increment(index int) VectorClock {
	out := VectorClock{}
	out.Values = append(out.Values, c.Values...)
	out.Values[index] += 1
	return out
}

func (c VectorClock) merge(other VectorClock) VectorClock {
	out := VectorClock{}
	out.Values = make([]int, len(c.Values))
	for i, v := range c.Values {
		if other.Values[i] > v {
			v = other.Values[i]
		}
		out.Values[i] = v
	}
	return out
}

func (c VectorClock) concurrentWith(other VectorClock) bool {
	return !c.happensBefore(other) && !other.happensBefore(c)
}

func (c VectorClock) String() string {
	return fmt.Sprint(c.Values)
}

func VectorClockWithValues(values ...int) VectorClock {
	return VectorClock{values}
}

func VectorClockWithZeros(number int) VectorClock {
	return VectorClock{make([]int, number)}
}

func VectorClockInfinity(number int) VectorClock {
	// computes the maximum int value (either 64 or 32 bit)
	var maxint = int(^uint(0) >> 1)
	values := make([]int, number)
	for i := range values {
		values[i] = maxint
	}
	return VectorClock{values}
}

var lineRegexp = regexp.MustCompile(`^([^0-9]*)(.*)$`)
var numberRegexp = regexp.MustCompile(`[^0-9]+`)

func Parse(input io.Reader) ([]VectorClock, map[string]string, error) {
	clocks := []VectorClock{}
	labels := map[string]string{}

	scanner := bufio.NewScanner(input)
	clockLength := 0
	for scanner.Scan() {
		// skip blank lines
		trimmed := strings.TrimSpace(scanner.Text())
		if trimmed == "" {
			continue
		}

		matches := lineRegexp.FindStringSubmatch(trimmed)
		if len(matches) == 0 {
			return nil, nil, fmt.Errorf("vclock.Parse: line %d does not appear to be valid", len(clocks)+1)
		}

		label := strings.TrimSpace(matches[1])
		numberParts := numberRegexp.Split(matches[2], -1)
		values := []int{}
		for _, n := range numberParts {
			if n == "" {
				continue
			}
			v, err := strconv.Atoi(n)
			if err != nil {
				return nil, nil, fmt.Errorf("vclock.Parse: line %d: %s", len(clocks)+1, err.Error())
			}
			values = append(values, v)
		}
		if len(values) == 0 {
			return nil, nil, fmt.Errorf("vclock.Parse: line %d no numbers", len(clocks)+1)
		}
		if clockLength > 0 && len(values) != clockLength {
			return nil, nil, fmt.Errorf("vclock.Parse: line %d has %d numbers; previous have %d", len(clocks)+1, len(values), clockLength)
		}
		clockLength = len(values)
		clock := VectorClockWithValues(values...)
		assert(len(clock.Values) == clockLength)
		clocks = append(clocks, clock)
		if label != "" {
			labels[clock.String()] = label
		}
	}
	if scanner.Err() != nil {
		return nil, nil, scanner.Err()
	}
	return clocks, labels, nil
}
