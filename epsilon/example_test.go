package epsilon

func ExampleNewEpsilonGreedy() {
	hp := NewEpsilonGreedy([]string{"a", "b"}, 0, &LinearEpsilonValueCalculator{})
	hostResponse := hp.Get()
	hostname := hostResponse.Host()
	hostResponse.Mark(nil)
	println(hostname)
}
