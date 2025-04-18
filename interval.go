package main

type interval struct {
	min, max float64
}

// func (i interval) size() float64 {
// 	return i.max - i.min
// }

// func (i interval) contains(x float64) bool {
// 	return i.min <= x && x <= i.max
// }

func (i interval) surrounds(x float64) bool {
	return i.min < x && x < i.max
}

func (i interval) clamp(x float64) float64 {
	if x < i.min {
		return i.min
	} else if x > i.max {
		return i.max
	}
	return x
}

// var empty interval = interval{min: math.Inf(1), max: math.Inf(-1)}
// var universe interval = interval{min: math.Inf(-1), max: math.Inf(1)}
