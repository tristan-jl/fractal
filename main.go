package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"time"
)

const (
	imgWidth               = 2560
	imgHeight              = 1440
	samplePerPixel         = 10
	centreRe       float64 = -0.5
	centreIm       float64 = 0.
	h              float64 = 2.4
	maxIter                = 25000
	escapeRadius   float64 = 1 << 8
	log2                   = 0.6931471805599453
)

const (
	pixelSize     = h / imgHeight // square pixels
	escapeRadius2 = escapeRadius * escapeRadius
)

func fractionalEscapeValue(cRe, cIm float64) float64 {
	x, y, xTemp := 0., 0., 0.

	for i := 0; i < maxIter; i++ {
		if x*x+y*y > escapeRadius2 {
			return float64(i) + 1. - math.Log(math.Log(x*x+y*y)/log2)/log2
		}
		xTemp = x*x - y*y + cRe
		y = 2*x*y + cIm
		x = xTemp
	}
	return float64(maxIter)
}

func setPixel(img *image.NRGBA, i, j int, cRe, cIm float64) {
	var value float64

	cReRands := randFloat(cRe, cRe+pixelSize, samplePerPixel)
	cImRands := randFloat(cIm, cIm+pixelSize, samplePerPixel)

	for i := 0; i < samplePerPixel; i++ {
		value += fractionalEscapeValue(cReRands[i], cImRands[i])
	}

	value = value / samplePerPixel

	if value < maxIter {
		r, g, b := hslToRgb(value/400., 1, 0.5)
		img.SetNRGBA(i, j, color.NRGBA{uint8(r * 255), uint8(g * 255), uint8(b * 255), 255})
	} else {
		img.SetNRGBA(i, j, color.NRGBA{0, 0, 0, 255})
	}
}

func main() {
	start := time.Now()
	fmt.Println("Starting")
	img := image.NewNRGBA(image.Rect(0, 0, imgWidth, imgHeight))

	for j := 0; j < imgHeight; j++ {
		cIm := centreIm + h/2 - float64(j)*pixelSize
		for i := 0; i < imgWidth; i++ {
			cRe := centreRe - h/2*(float64(imgWidth)/float64(imgHeight)) + float64(i)*pixelSize

			go setPixel(img, i, j, cRe, cIm)
		}
		fmt.Printf("\r%d/%d", j+1, imgHeight) // TODO make better progress bar
	}
	fmt.Println()

	f, err := os.Create("result.png")
	if err != nil {
		panic(err)
	}
	err = png.Encode(f, img)
	if err != nil {
		panic(err)
	}
	elapsed := time.Since(start)
	fmt.Printf("Done in %s\n", elapsed)
}
