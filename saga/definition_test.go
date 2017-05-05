package saga

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func T1(ctx context.Context) {

}

func C1(ctx context.Context) {

}

func T2(ctx context.Context) {

}

func C2(ctx context.Context) {

}

func TestSubTxDefinitions(t *testing.T) {
	txs := subTxDefinitions{}.
		addDefinition("A1", T1, C1).
		addDefinition("A2", T2, C2)
	define, ok := txs.findDefinition("A1")
	assert.True(t, ok)
	assert.NotNil(t, define.action)
}

func E() {

}
