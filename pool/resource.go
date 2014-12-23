package pool

// Factory is a function that can be used to create a resource.
type Factory func() (Resource, error)

// Every resource needs to suport the Resource interface.
// A ResourcePool is composed of several Resource's
type Resource interface {
	Id() uint64
	Close()
}
