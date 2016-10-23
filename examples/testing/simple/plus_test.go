package simple

import (
	"testing"
)

func TestPlusShouldSucceed(t *testing.T) {
	a := 1
	b := 1
	shouldBe := 2
	if real := plus(a, b); real != shouldBe {
		t.Errorf("plus(%d, %d) should be %d, but is:%d\n", a, b, shouldBe, real)
	}
}
