package pool

// Factory is a function that can be used to create a resource.
type Factory func() (Resource, error)

// Every resource needs to suport the Resource interface.
// Thread synchronization between Close() and IsClosed()
// is the responsibility the caller.
type Resource interface {
	Close()
}
