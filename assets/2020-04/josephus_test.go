package josephus

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJosephusBitMap(t *testing.T) {
	testcases := map[int]int{
		1:  1,
		2:  1,
		3:  3,
		10: 5,
	}

	for input, expect := range testcases {
		assert.Equal(t, expect, JosephusBitMap(input), fmt.Sprintf("input=%d", input))
	}
}

func BenchmarkJosephusBitMap10(b *testing.B) {
	for i := 0; i < b.N; i++ {
		JosephusBitMap(10)
	}
}

func BenchmarkJosephusBitMap10000(b *testing.B) {
	for i := 0; i < b.N; i++ {
		JosephusBitMap(10000)
	}
}

func TestJosephusLinklist(t *testing.T) {
	testcases := map[int]int{
		1:  1,
		2:  1,
		3:  3,
		10: 5,
	}

	for input, expect := range testcases {
		assert.Equal(t, expect, JosephusLinklist(input), fmt.Sprintf("input=%d", input))
	}
}

func BenchmarkJosephusLinklist10(b *testing.B) {
	for i := 0; i < b.N; i++ {
		JosephusLinklist(10)
	}
}

func BenchmarkJosephusLinklist10000(b *testing.B) {
	for i := 0; i < b.N; i++ {
		JosephusLinklist(10000)
	}
}

func BenchmarkJosephusRecursion10(b *testing.B) {
	for i := 0; i < b.N; i++ {
		JosephusRecursion(10)
	}
}

func BenchmarkJosephusRecursion10000(b *testing.B) {
	for i := 0; i < b.N; i++ {
		JosephusRecursion(10000)
	}
}

func TestJosephusTable(t *testing.T) {
	fmt.Printf(" n,f(n), i| %5s, %5s, %5s\n", "n_bit", "fn_bi", "i_bit")
	for i := 1; i <= 20; i++ {
		res := JosephusRecursion(i)
		// fmt.Printf("%2d, %2d\n", i, res)
		fmt.Printf("%2d, %2d, %2d| %05b, %05b, %05b\n", i, res, res/2, i, res, res/2)
	}
}

func TestJosephusBit(t *testing.T) {
	testcases := map[int]int{
		1:  1,
		2:  1,
		3:  3,
		10: 5,
	}

	for input, expect := range testcases {
		assert.Equal(t, expect, JosephusBit(input), fmt.Sprintf("input=%d", input))
	}
}

func BenchmarkJosephusBit10(b *testing.B) {
	for i := 0; i < b.N; i++ {
		JosephusBit(10)
	}
}

func BenchmarkJosephusBit10000(b *testing.B) {
	for i := 0; i < b.N; i++ {
		JosephusBit(10000)
	}
}
