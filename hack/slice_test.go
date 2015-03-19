package hack

import (
	"github.com/funkygao/assert"
	"testing"
)

func TestAppend(t *testing.T) {
	a := make([]interface{}, 3, 3)
	a[0] = 1
	a[1] = 5
	a[2] = 9
	t.Logf("%+v", a)
	assert.Equal(t, 3, len(a))
	assert.Equal(t, 3, cap(a))
	c := append(a, 11)
	assert.Equal(t, 4, len(c))
	assert.Equal(t, 2*cap(a), cap(c)) // internally slice grow by factor of 2

	b := Append(a, 11, 2) // grow by 2
	t.Logf("%+v", b)
	assert.Equal(t, []interface{}{1, 5, 9, 11}, b)
	assert.Equal(t, 3+2, cap(b))

}
