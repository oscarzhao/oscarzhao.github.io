package simple

import (
	"errors"
	"strconv"
	"strings"
)

var (
	errInvalidIPAddrMissingParts = errors.New("invalid ip address: missing parts")
	errInvalidIPAddrNotNumber    = errors.New("invalid ip address: not valid number")
)

// ParseIP consumes a IP address string, and returns corresponding ip array
func ParseIP(str string) ([]int, error) {
	items := strings.Split(str, ".")
	if len(items) < 4 {
		return nil, errInvalidIPAddrMissingParts
	}

	var res []int
	for _, item := range items {
		val, err := strconv.ParseInt(item, 10, 32)
		if err != nil || val < 0 || val > 255 {
			return nil, errInvalidIPAddrNotNumber
		}

		res = append(res, int(val))
	}

	return res, nil
}
