package simple

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ParseIP(t *testing.T) {
	testcases := []struct {
		desc      string
		input     string
		expect    []int
		expectErr error
	}{
		{
			desc:   "valid ip 001",
			input:  "0.0.0.0",
			expect: []int{0, 0, 0, 0},
		},
		{
			desc:   "valid ip 002",
			input:  "255.255.255.255",
			expect: []int{255, 255, 255, 255},
		},
		{
			desc:      "too big number",
			input:     "1234.11.0.0",
			expectErr: errInvalidIPAddrNotNumber,
		},
		{
			desc:      "invalid character",
			input:     "12ab.11.0.0",
			expectErr: errInvalidIPAddrNotNumber,
		},
	}

	for _, ts := range testcases {
		got, gotErr := ParseIP(ts.input)
		require.Equal(t, ts.expect, got, ts.desc)
		require.Equal(t, ts.expectErr, gotErr, ts.desc)
	}
}
