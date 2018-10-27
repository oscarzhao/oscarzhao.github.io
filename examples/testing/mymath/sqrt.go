package mymath

// Sqrt calculate the square root of a positive float64 number
// max error is about 10^-9. For simplicity, we didn't discard
// the part with is smaller than 10^-9
func Sqrt(x float64) float64 {
	if x < 0 {
		panic("cannot be negative")
	}

	if x == 0 {
		return 0
	}

	a := x / 2
	b := (a + 2) / 2
	err := a - b
	for err >= 0.000000001 || err <= -0.000000001 {
		a = b
		b = (b + x/b) / 2
		err = a - b
	}

	return b
}
