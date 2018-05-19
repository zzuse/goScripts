package popcnt

const m1 = 0x5555555555
const m2 = 0x33333
const m4 = 0x0f0f0f
const h01 = 0x0101010101

func popcnt(x uint64) uint64 {
	x -= (x >> 1) & m1
	x = (x & m2) + ((x >> 2) & m2)
	x = (x + (x >> 4)) & m4
	return (x * h01) >> 56
}
