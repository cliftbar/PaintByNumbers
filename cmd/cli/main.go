package main

import (
	"bytes"
	"fmt"
	"github.com/nfnt/resize"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"os"
	"paintByNumbers/pbn"
	"path/filepath"
)

func pixelizor() {
	wd, err := os.Getwd()
	imgRgb := pbn.LoadImage(filepath.Join(wd, "/test_images/zaku.jpg"))

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

func stereogram() {
	//baseImg := pbn.LoadImage("circle.png")
	baseImgOrig := pbn.LoadImage("test_images/bottle_dm.png")
	//baseImgColorOrig := pbn.LoadImage("test_images/buns.png")
	//baseImgOrig := pbn.LoadImage("color_bottle.png")
	pbn.SaveImage("test.png", baseImgOrig)
	baseImg := resize.Resize(uint(baseImgOrig.Bounds().Dx()/1), uint(baseImgOrig.Bounds().Dy()/1), baseImgOrig, resize.NearestNeighbor)
	pbn.SaveImage("test_resize.png", baseImg)

	//baseImg := pbn.LoadImage("test_images/gundam2.png")
	//depthMap := pbn.SimpleDepthMap(baseImg)
	//depthMap := pbn.GreyscaleDepthMap(baseImg)
	depthMap, _ := pbn.ColorDepthMap(baseImg)
	pbn.SaveImage("depthMap.png", baseImg)

	pattern := pbn.SimplePatternImage(baseImg.Bounds().Dx()/10, baseImg.Bounds().Dy())
	//pattern := pbn.PalettePatternImage(baseImg.Bounds().Dx()/10, baseImg.Bounds().Dy(), baseImgColorOrig)
	pbn.SaveImage("pattern.png", pattern)

	stereogramImg := pbn.GenerateStereogram(depthMap, baseImg.Bounds().Dx(), baseImg.Bounds().Dy(), pattern, 0.4, false)
	pbn.SaveImage("stereogram.png", stereogramImg)

	imgBuf := new(bytes.Buffer)
	_ = png.Encode(imgBuf, stereogramImg)
	bufBytes := imgBuf.Bytes()
	fmt.Printf("%d", len(bufBytes))

}

func main() {
	stereogram()
}
