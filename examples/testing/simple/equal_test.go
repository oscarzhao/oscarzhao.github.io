package simple

import (
	"testing"
)

func TestEqualShouldFail(t *testing.T) {
	a := 1
	b := 1
	shouldNotBe := false
	if real := equal(a, b); real == shouldNotBe {
		t.Errorf("equal(%d, %d) should not be %v, but is:%v\n", a, b, shouldNotBe, real)
	}
}
