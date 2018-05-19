package fib

import "testing"

func BenchmarkFib(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Fib(20)
	}
	// Output: MOOOO!
}
