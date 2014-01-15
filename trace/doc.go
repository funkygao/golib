/*
defer is LIFO

Usage:

    import t "kx/trace"

    func foo() {
        defer t.Un(t.Trace("foo"))

        //
    }
*/
package trace
