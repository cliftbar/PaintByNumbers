package main

import (
	"fmt"
	"github.com/nfnt/resize"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"os"
	"paintByNumbers/pbn"
	"path/filepath"
)

func main() {
	wd, err := os.Getwd()
	imgRgb := pbn.LoadImage(filepath.Join(wd, "/test_images/pokemon_mindjourney.png"))

	//resizedImgRgb := image.NewRGBA(image.Rect(0, 0, imgRgb.Bounds().Size().X/10, imgRgb.Bounds().Size().Y/10))
	//err = rez.Convert(resizedImgRgb, imgRgb, rez.NewBilinearFilter())
	widthXFactor := 5
	heightYFactor := 5
	resizedImgRgb := resize.Resize(uint(imgRgb.Bounds().Size().X/widthXFactor), uint(imgRgb.Bounds().Size().Y/heightYFactor), imgRgb, resize.NearestNeighbor)

	fmt.Println("image loaded")
	colorPalette := pbn.DominantColors(imgRgb, 5, 0.01, true)
	fmt.Println("color palette found")
	snapImg := pbn.SnapColors(resizedImgRgb, colorPalette)
	fmt.Println("snap done")

	//bigImgRgb := image.NewRGBA(image.Rect(0, 0, imgRgb.Bounds().Size().X, imgRgb.Bounds().Size().Y))
	//err = rez.Convert(bigImgRgb, snapImg, rez.NewBilinearFilter())

	newImage := resize.Resize(uint(imgRgb.Bounds().Size().X), uint(imgRgb.Bounds().Size().Y), snapImg, resize.NearestNeighbor)

	outputFile, err := os.Create("out.png")
	if err != nil {
		// Handle error
	}

	err = png.Encode(outputFile, newImage)
	if err != nil {
		// Handle error
	}
}
