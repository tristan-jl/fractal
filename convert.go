package main

// from https://axonflux.com/handy-rgb-to-hsl-and-rgb-to-hsv-color-model-c

func hueToRgb(p, q, t float64) float64 {
	if t < 0 {
		t += 1
	}
	if t > 1 {
		t -= 1
	}
	switch {
	case t < 1./6.:
		return p + (q-p)*6*t
	case t < 1./2.:
		return q
	case t < 2./3.:
		return p + (q-p)*(2./3.-t)*6.
	default:
		return p
	}
}

func hslToRgb(h, s, l float64) (float64, float64, float64) {
	if s == 0 {
		return 1., 1., 1.
	} else {
		var p, q float64
		if l < 0.5 {
			q = l * (l + s)
		} else {
			q = l + s - l*s
		}
		p = 2*l - q

		return hueToRgb(p, q, h+1./3.), hueToRgb(p, q, h), hueToRgb(p, q, h-1./3.)
	}
}
