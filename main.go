package main

import (
	"bytes"
	"fmt"
	clf "github.com/lucasb-eyer/go-colorful"
	"github.com/nfnt/resize"
	_ "golang.org/x/image/webp"
	"image"
	_ "image/jpeg"
	"image/png"
	"paintByNumbers/pbn"
	"strings"
	"syscall/js"
	"time"
)

func pixelizor(this js.Value, i []js.Value) interface{} {
	start := time.Now()

	height := i[0].Int()
	width := i[1].Int()
	widthXScalar := i[2].Int()
	heightYScalar := i[3].Int()
	clusterCount := i[4].Int()
	deltaThreshold := i[5].Float()
	srcArrayJS := i[6]

	outputBuffer := i[7]

	srcLen := srcArrayJS.Get("byteLength").Int()
	inputImageBytes := make([]uint8, srcLen)
	js.CopyBytesToGo(inputImageBytes, srcArrayJS)
	fmt.Printf("src bytes copied %f\n", time.Since(start).Seconds())

	imgRgb, inputType, err := image.Decode(bytes.NewReader(inputImageBytes))
	fmt.Printf("decode image done %f\n", time.Since(start).Seconds())
	if err != nil {
		panic("unable to read image")
	}
	fmt.Printf("Height: %d, width: %d, pixels: %d, type: %s\n", height, width, len(inputImageBytes), inputType)

	resizedImgRgb := resize.Resize(uint(imgRgb.Bounds().Size().X/widthXScalar), uint(imgRgb.Bounds().Size().Y/heightYScalar), imgRgb, resize.NearestNeighbor)
	fmt.Printf("image resize down %f\n", time.Since(start).Seconds())

	colorPalette := pbn.DominantColors(imgRgb, clusterCount, deltaThreshold, false)
	fmt.Printf("color palette found %f\n", time.Since(start).Seconds())

	snapImg := pbn.SnapColors(resizedImgRgb, colorPalette)
	fmt.Printf("snap done %f\n", time.Since(start).Seconds())

	newImage := resize.Resize(uint(imgRgb.Bounds().Size().X), uint(imgRgb.Bounds().Size().Y), snapImg, resize.NearestNeighbor)
	fmt.Printf("resize up done %f\n", time.Since(start).Seconds())

	var imgBuf *bytes.Buffer
	var bufBytes []uint8

	imgBuf = new(bytes.Buffer)
	_ = png.Encode(imgBuf, newImage)
	bufBytes = imgBuf.Bytes()

	js.CopyBytesToJS(outputBuffer, bufBytes)
	fmt.Printf("outputBuffer bytes copied %f\n", time.Since(start).Seconds())

	colorPaletteHex := make([]string, 0)
	colorPaletteHexStr := ""
	for _, c := range colorPalette {
		hex := c.Hex()
		colorPaletteHex = append(colorPaletteHex, hex)
		if colorPaletteHexStr == "" {
			colorPaletteHexStr = hex
		} else {
			colorPaletteHexStr = colorPaletteHexStr + "," + hex
		}
	}

	fmt.Printf("Wasm Done %fs\n", time.Since(start).Seconds())
	return js.ValueOf(colorPaletteHexStr)
}

func dominantColors(this js.Value, i []js.Value) interface{} {
	start := time.Now()

	height := i[0].Int()
	width := i[1].Int()
	clusterCount := i[2].Int()
	deltaThreshold := i[3].Float()
	srcArrayJS := i[4]

	srcLen := srcArrayJS.Get("byteLength").Int()
	inputImageBytes := make([]uint8, srcLen)
	js.CopyBytesToGo(inputImageBytes, srcArrayJS)
	fmt.Printf("src bytes copied %f\n", time.Since(start).Seconds())

	imgRgb, inputType, err := image.Decode(bytes.NewReader(inputImageBytes))
	fmt.Printf("decode image done %f\n", time.Since(start).Seconds())
	if err != nil {
		panic("unable to read image")
	}
	fmt.Printf("Height: %d, width: %d, pixels: %d, type: %s\n", height, width, len(inputImageBytes), inputType)

	colorPalette := pbn.DominantColors(imgRgb, clusterCount, deltaThreshold, false)
	fmt.Printf("color palette found, %d clusters, %f threshold %f\n", clusterCount, deltaThreshold, time.Since(start).Seconds())

	colorPaletteHex := make([]string, 0)
	colorPaletteHexStr := ""
	for _, c := range colorPalette {
		hex := c.Hex()
		colorPaletteHex = append(colorPaletteHex, hex)
		if colorPaletteHexStr == "" {
			colorPaletteHexStr = hex
		} else {
			colorPaletteHexStr = colorPaletteHexStr + "," + hex
		}
	}

	fmt.Printf("Wasm Done %fs\n", time.Since(start).Seconds())
	return js.ValueOf(colorPaletteHexStr)
}

func pixelizeFromPalette(this js.Value, i []js.Value) interface{} {
	start := time.Now()

	height := i[0].Int()
	width := i[1].Int()
	widthXScalar := i[2].Int()
	heightYScalar := i[3].Int()
	colorsHexString := i[4].String()
	srcArrayJS := i[5]

	outputBuffer := i[6]

	colorsHex := strings.Split(colorsHexString, ",")
	colorPalette := make([]clf.Color, 0)
	for _, cHex := range colorsHex {
		c, _ := clf.Hex(cHex)
		colorPalette = append(colorPalette, c)
	}

	srcLen := srcArrayJS.Get("byteLength").Int()
	inputImageBytes := make([]uint8, srcLen)
	js.CopyBytesToGo(inputImageBytes, srcArrayJS)
	fmt.Printf("src bytes copied %f\n", time.Since(start).Seconds())

	imgRgb, inputType, err := image.Decode(bytes.NewReader(inputImageBytes))
	fmt.Printf("decode image done %f\n", time.Since(start).Seconds())
	if err != nil {
		panic("unable to read image")
	}
	fmt.Printf("Height: %d, width: %d, pixels: %d, type: %s\n", height, width, len(inputImageBytes), inputType)

	resizedImgRgb := resize.Resize(uint(imgRgb.Bounds().Size().X/widthXScalar), uint(imgRgb.Bounds().Size().Y/heightYScalar), imgRgb, resize.NearestNeighbor)
	fmt.Printf("image resize down %f\n", time.Since(start).Seconds())

	snapImg := pbn.SnapColors(resizedImgRgb, colorPalette)
	fmt.Printf("snap done %f\n", time.Since(start).Seconds())

	newImage := resize.Resize(uint(imgRgb.Bounds().Size().X), uint(imgRgb.Bounds().Size().Y), snapImg, resize.NearestNeighbor)
	fmt.Printf("resize up done %f\n", time.Since(start).Seconds())

	var imgBuf *bytes.Buffer
	var bufBytes []uint8

	imgBuf = new(bytes.Buffer)
	_ = png.Encode(imgBuf, newImage)
	bufBytes = imgBuf.Bytes()

	js.CopyBytesToJS(outputBuffer, bufBytes)
	fmt.Printf("outputBuffer bytes copied %f\n", time.Since(start).Seconds())

	fmt.Printf("Wasm Done %fs\n", time.Since(start).Seconds())
	return js.ValueOf("done")
}

func registerCallbacks() {
	js.Global().Set("pixelizor", js.FuncOf(pixelizor))
	js.Global().Set("dominantColors", js.FuncOf(dominantColors))
	js.Global().Set("pixelizeFromPalette", js.FuncOf(pixelizeFromPalette))
}

func main() {
	c := make(chan struct{}, 0)

	println("WASM Go Initialized")
	// register functions
	registerCallbacks()
	<-c
}
