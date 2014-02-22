package pool

import (
	"fmt"
	"time"
)

func (this *ResourcePool) StatsJSON() string {
	if this == nil {
		return "{}"
	}
	c, a, mx, wc, wt, it := this.Stats()
	return fmt.Sprintf(`{"Capacity": %v, "Available": %v, "MaxCapacity": %v, "WaitCount": %v, "WaitTime": %v, "IdleTimeout": %v}`,
		c, a, mx, wc, int64(wt), int64(it))
}

func (this *ResourcePool) Stats() (capacity, available, maxCap, waitCount int64, waitTime, idleTimeout time.Duration) {
	return this.Capacity(), this.Available(),
		this.MaxCapacity(), this.WaitCount(),
		this.WaitTime(), this.IdleTimeout()
}
