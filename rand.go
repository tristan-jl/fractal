package main

import "math/rand"

func randFloat(min, max float64, n int) []float64 {
	arr := make([]float64, n)
	for i := range arr {
		arr[i] = min + rand.Float64()*(max-min)
	}
	return arr
}
