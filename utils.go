package main

import (
	"math"
	"math/rand/v2"
)

func deg2rad(degrees float64) float64 {
	return degrees * math.Pi / 180.0
}

func random() float64 {
	return rand.Float64()
}

func randomIn(min, max float64) float64 {
	return min + (max-min)*random()
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
