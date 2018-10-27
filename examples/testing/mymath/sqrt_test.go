package mymath

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSqrt(t *testing.T) {
	testcases := []struct {
		desc   string
		input  float64
		expect float64
	}{
		{
			desc:   "zero",
			input:  0,
			expect: 0,
		},
		{
			desc:   "one",
			input:  1,
			expect: 1,
		},
		{
			desc:   "a very small rational number",
			input:  0.0000000000000000000001,
			expect: 0.0,
		},
		{
			desc:   "rational number result: 2.56",
			input:  2.56,
			expect: 1.6,
		},
		{
			desc:   "irrational number result: 2",
			input:  2,
			expect: 1.414213562,
		},
	}

	for _, ts := range testcases {
		got := Sqrt(ts.input)
		err := got - ts.expect
		require.True(t, err < 0.000000001 && err > -0.00000001, ts.desc)
	}
}

func TestSqrt_Panic(t *testing.T) {
	defer func() {
		r := recover()
		require.Equal(t, "cannot be negative", r)
	}()
	_ = Sqrt(-1)
}

func Test_Require_EqualValues(t *testing.T) {
	// tests will pass
	require.EqualValues(t, 12, 12.0, "compare int32 and float64")
	require.EqualValues(t, 12, int64(12), "compare int32 and int64")

	// tests will fail
	require.Equal(t, 12, 12.0, "compare int32 and float64")
	require.Equal(t, 12, int64(12), "compare int32 and int64")
}
