package saga

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshalLog(t *testing.T) {
	l := &Log{
		Type:    ActionStart,
		SubTxID: "1",
		Params:  []ParamData{},
	}
	sl := l.mustMarshal()
	l2 := mustUnmarshalLog(sl)
	assert.Equal(t, ActionStart, l2.Type)
}
