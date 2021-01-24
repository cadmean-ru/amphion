package native

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type testFeature struct {
	result string
}

func (f *testFeature) OnWeb() {
	f.result = "test web"
}

func (f *testFeature) OnPc() {
	f.result = "test pc"
}

func TestInvoke(t *testing.T) {
	var f = &testFeature{}
	Invoke(f)
	var expected = "test pc"
	assert.Equal(t, expected, f.result)
}