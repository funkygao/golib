package pool

// Factory is a function that can be used to create a resource.
type Factory func() (Resource, error)

// Every resource needs to suport the Resource interface.
type Resource interface {
	Close()
}

// For those resource that wants keepalive feature, it will
// implement Keepalive interface.
type Keepalive interface {
	Keepalive()
}
