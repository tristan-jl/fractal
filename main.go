package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
)

const (
	imgWidth               = 2560
	imgHeight              = 1440
	samplePerPixel         = 1
	minRe          float64 = -2.5
	minIm          float64 = -1.2
	maxIm          float64 = 1.2
	maxIter                = 25
	escapeRadius   float64 = 2
)

const (
	maxRe         = minRe + (maxIm-minIm)*imgWidth/imgHeight
	rePixelSize   = (maxRe - minRe) / (imgWidth - 1)
	imPixelSize   = (maxIm - minIm) / (imgHeight - 1)
	escapeRadius2 = escapeRadius * escapeRadius
)

func escapeValue(cRe, cIm float64) (float64, int) {
	x, y, x2, y2 := 0., 0., 0., 0.

	for i := 0; i < maxIter; i++ {
		if x2+y2 > escapeRadius2 {
			return x2 + y2, i
		}
		y = 2*x*y + cIm
		x = x2 - y2 + cRe
		x2 = x * x
		y2 = y * y
	}

	return x2 + y2, maxIter
}

func fractionalEscapeValue(cRe, cIm float64) float64 {
	z, n := escapeValue(cRe, cIm)

	logBase := 1. / math.Log(2.)
	logHalfBase := math.Log(0.5) * logBase

	return 5. + float64(n) - logHalfBase - math.Log(math.Log(z))*logBase
}

func main() {
	fmt.Println("Starting")
	img := image.NewNRGBA(image.Rect(0, 0, imgWidth, imgHeight))

	for j := 0; j < imgHeight; j++ {
		cIm := maxIm - float64(j)*imPixelSize
		for i := 0; i < imgWidth; i++ {
			cRe := minRe + float64(i)*rePixelSize

			value := fractionalEscapeValue(cRe, cIm)

			if value < maxIter {
				r, g, b := hslToRgb(value/float64(maxIter), 1, 0.5)
				img.SetNRGBA(i, j, color.NRGBA{uint8(r * 255), uint8(g * 255), uint8(b * 255), 255})
			} else {
				img.SetNRGBA(i, j, color.NRGBA{0, 0, 0, 255})
			}
		}
	}

	f, err := os.Create("result.png")
	if err != nil {
		panic(err)
	}
	err = png.Encode(f, img)
	if err != nil {
		panic(err)
	}
	fmt.Println("Done!")
}
