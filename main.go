package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
)

const (
	imgWidth               = 2560
	imgHeight              = 1440
	samplePerPixel         = 1
	maxIter                = 100
	minRe          float64 = -2.5
	minIm          float64 = -1.2
	maxIm          float64 = 1.2
)

const (
	maxRe    = minRe + (maxIm-minIm)*imgWidth/imgHeight
	reFactor = (maxRe - minRe) / (imgWidth - 1)
	imFactor = (maxIm - minIm) / (imgHeight - 1)
)

func escapeTime(cRe, cIm float64) int {
	var x, y, x2, y2 float64

	for i := 0; i < maxIter; i++ {
		if x2+y2 > 4 {
			return i
		}
		y = 2*x*y + cIm
		x = x2 - y2 + cRe
		x2 = x * x
		y2 = y * y
	}

	return maxIter
}

func main() {
	fmt.Println("Starting")
	img := image.NewNRGBA(image.Rect(0, 0, imgWidth, imgHeight))

	for j := 0; j < imgHeight; j++ {
		cIm := maxIm - float64(j)*imFactor
		for i := 0; i < imgWidth; i++ {
			cRe := minRe + float64(i)*reFactor

			iter := escapeTime(cRe, cIm)

			if iter >= maxIter {
				img.SetNRGBA(i, j, color.NRGBA{0, 0, 0, 255})
			} else {
				r, g, b := hslToRgb(float64(iter)/float64(maxIter), 1, 0.5)
				img.SetNRGBA(i, j, color.NRGBA{uint8(r * 255), uint8(g * 255), uint8(b * 255), 255})
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
