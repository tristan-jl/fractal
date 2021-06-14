package main

import (
	"image/color"
	"math"
)

// from https://en.wikipedia.org/wiki/Monotone_cubic_interpolation

// xs, ys must be sorted
func createInterpolant(xs, ys []float64) func(float64) float64 {
	length := len(xs)

	if length != len(ys) {
		panic("Length of arguments must be the same.")
	}
	if length == 0 {
		return func(x float64) float64 { return 0. }
	}
	if length == 1 {
		return func(x float64) float64 { return ys[0] }
	}

	var dys, dxs, ms []float64
	for i := 0; i < length-1; i++ {
		dx := xs[i+1] - xs[i]
		dy := ys[i+1] - ys[i]
		dxs = append(dxs, dx)
		dys = append(dys, dy)
		ms = append(ms, dy/dx)
	}

	cls := []float64{ms[0]}
	for i := 0; i < len(dxs)-1; i++ {
		m := ms[i]
		mNext := ms[i+1]

		if m*mNext <= 0 {
			cls = append(cls, 0)
		} else {
			dx_ := dxs[i]
			dxNext := dxs[i+1]
			common := dx_ + dxNext
			cls = append(cls, 3*common/((common+dxNext)/m+(common+dx_)/mNext))
		}
	}
	cls = append(cls, ms[len(ms)-1])

	var c2s, c3s []float64
	for i := 0; i < len(cls)-1; i++ {
		c1 := cls[i]
		m_ := ms[i]
		invDx := 1 / dxs[i]
		common_ := c1 + cls[i+1] - m_ - m_

		c2s = append(c2s, (m_-c1-common_)*invDx)
		c3s = append(c3s, common_*invDx*invDx)
	}

	return func(x float64) float64 {
		i := len(xs) - 1
		if x == xs[i] {
			return ys[i]
		}

		var low, mid int
		high := len(c3s) - 1

		for low <= high {
			mid = int(math.Floor(0.5 * (float64(low) + float64(high))))
			xHere := xs[mid]
			if xHere < x {
				low = mid + 1
			} else if xHere > x {
				high = mid - 1
			} else {
				return ys[mid]
			}
		}
		j := 0
		if high > 0 {
			j = high
		}

		diff := x - xs[j]
		diffSq := diff * diff
		return ys[j] + cls[j]*diff + c2s[j]*diffSq + c3s[j]*diff*diffSq
	}
}

func gradient(x float64) color.NRGBA {
	positions := []float64{0., 0.16, 0.42, 0.6425, 0.8575}

	red := createInterpolant(positions, []float64{0., 32., 237., 255., 0.})(x)
	green := createInterpolant(positions, []float64{7., 107., 255., 170., 2.})(x)
	blue := createInterpolant(positions, []float64{100., 203., 255., 0., 0.})(x)

	return color.NRGBA{uint8(red), uint8(green), uint8(blue), 255}
}
