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
	samplePerPixel         = 100
	centreRe       float64 = -0.5
	centreIm       float64 = 0.
	h              float64 = 2.4
	maxIter                = 25000
	escapeRadius   float64 = 1 << 8
)

const (
	gradientScale = 128
	numColours    = 2048
)

const (
	pixelSize     = h / imgHeight // square pixels
	escapeRadius2 = escapeRadius * escapeRadius
	log2          = 0.6931471805599453
)

func fractionalEscapeValue(cRe, cIm float64) (float64, int) {
	x, y, xTemp := 0., 0., 0.

	for i := 0; i < maxIter; i++ {
		if x*x+y*y > escapeRadius2 {
			return float64(i) + 10. - math.Log(math.Log(x*x+y*y)/log2)/log2, i
		}
		xTemp = x*x - y*y + cRe
		y = 2*x*y + cIm
		x = xTemp
	}
	return float64(maxIter), maxIter
}

func setPixel(img *image.NRGBA, colourGradient *[numColours]color.NRGBA, i, j int, cRe, cIm float64) {
	var value, valueTemp, iterFloat float64
	var iter, iterTemp int

	cReRands := randFloat(cRe, cRe+pixelSize, samplePerPixel)
	cImRands := randFloat(cIm, cIm+pixelSize, samplePerPixel)

	for i := 0; i < samplePerPixel; i++ {
		valueTemp, iterTemp = fractionalEscapeValue(cReRands[i], cImRands[i])
		value += valueTemp
		iter += iterTemp
	}

	value /= samplePerPixel
	iterFloat = float64(iter) / samplePerPixel

	if value < maxIter {
		colourI := math.Mod((math.Sqrt(iterFloat) * gradientScale), numColours)

		colour1 := colourGradient[int(colourI)]
		colour2 := colourGradient[(int(colourI)+1)%numColours]
		img.SetNRGBA(i, j, linearInterpolate(colour1, colour2, math.Mod(colourI, 1)))
	} else {
		img.SetNRGBA(i, j, color.NRGBA{0, 0, 0, 255})
	}
}

func main() {
	start := time.Now()
	fmt.Println("Starting")
	img := image.NewNRGBA(image.Rect(0, 0, imgWidth, imgHeight))

	var colourGradient [numColours]color.NRGBA
	for i := 0; i < numColours; i++ {
		colourGradient[i] = gradient(float64(i) / numColours)
	}

	for j := 0; j < imgHeight; j++ {
		cIm := centreIm + h/2 - float64(j)*pixelSize
		for i := 0; i < imgWidth; i++ {
			cRe := centreRe - h/2*(float64(imgWidth)/float64(imgHeight)) + float64(i)*pixelSize

			go setPixel(img, &colourGradient, i, j, cRe, cIm)
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
