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
	maxIter                = 30
	minRe          float64 = -2
	maxRe          float64 = 1
	minIm          float64 = -0.8
)

const (
	maxIm    = minIm + (maxRe-minRe)*imgHeight/imgWidth
	reFactor = (maxRe - minRe) / (imgWidth - 1)
	imFactor = (maxIm - minIm) / (imgHeight - 1)
)

func pixelColour(cRe, cIm float64) int {
	zRe := cRe
	zIm := cIm

	for i := 0; i < maxIter; i++ {
		zRe2, zReIm, zIm2 := zRe*zRe, zRe*zIm, zIm*zIm

		if zRe2+zIm2 > 4 {
			return i
		}
		zRe = zRe2 - zIm2 + cRe
		zIm = 2*zReIm + cIm
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

			pixel := pixelColour(cRe, cIm)

			if pixel >= maxIter {
				img.SetNRGBA(i, j, color.NRGBA{0, 0, 0, 255})
			} else {
				img.SetNRGBA(i, j, color.NRGBA{255, 255, 255, 255})
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
