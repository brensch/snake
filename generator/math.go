package generator

// Some int maths helpers to hopefully make things a bit faster
func Abs(num int32) int32 {
	if num < 0 {
		return -num
	}
	return num
}

func IntSqrt(n int32) int32 {

	x := n
	y := int32(1)
	for x > y {
		x = (x + y) / 2
		y = n / x
	}

	return x
}
