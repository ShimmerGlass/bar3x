package markup

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

type testEmbedT struct {
	testT
}

type testT struct {
	cb func(int)
	v  int
}

func (t *testT) Val() int {
	return t.v
}

func (t *testT) SetVal(v int) {
	t.v = v
	if t.cb != nil {
		t.cb(v)
	}
}

func (t *testT) OnValChange(cb func(int)) {
	t.cb = cb
	cb(t.v)
}

func TestField(t *testing.T) {
	v := &testEmbedT{}
	f, err := newField(reflect.ValueOf(v), "Val")
	require.Nil(t, err)

	require.Equal(t, 0, f.Get())
	require.Nil(t, f.Set(5))
	require.Equal(t, 5, f.Get())

	cbVals := []interface{}{}

	require.Nil(t, f.Watch(func(v interface{}) {
		cbVals = append(cbVals, v)
	}))

	require.Nil(t, f.Set(10))

	require.Equal(t, []interface{}{5, 10}, cbVals)
}
