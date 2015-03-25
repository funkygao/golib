package assert

import (
	"fmt"
)

func Assert(condition bool, msg string, v ...interface{}) {
	if !condition {
		panic(fmt.Sprintf("assert failed "+msg, v...))
	}
}
