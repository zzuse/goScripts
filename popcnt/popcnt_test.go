package popcnt

import "testing"

var result uint64

func benchmarkPopcnt(b *testing.B) {
	var r uint64
	for i := 0; i < b.N; i++ {
		r = popcnt(uint64(i))
	}
	result = r
}

//func benchmarkPopcnt(b *testing.B) {
//	for i := 0; i < b.N; i++ {
//		popcnt(uint64(i)) //optimied away
//	}
//}
