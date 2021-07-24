package gpu

import (
	"fmt"
	"testing"
)

func TestMemoryLayout_Stride(t *testing.T) {
	var ml = NewMemoryLayout(struct {
		name string
		test string
	}{})
	ml.Calculate()
	fmt.Println(ml.Stride())
}
