package pool

import (
	"fmt"
	"time"
)

func (this *ResourcePool) StatsJSON() string {
	if this == nil {
		return "{}"
	}
	c, a, mx, wc, it := this.Stats()
	return fmt.Sprintf(`{"Capacity": %v, "Available": %v, "MaxCapacity": %v, "WaitCount": %v, "IdleTimeout": %v}`,
		c, a, mx, wc, int64(it))
}

func (this *ResourcePool) Stats() (capacity, available, maxCap, waitCount int64, idleTimeout time.Duration) {
	return this.Capacity(), this.Available(),
		this.MaxCapacity(), this.WaitCount(),
		this.IdleTimeout()
}
