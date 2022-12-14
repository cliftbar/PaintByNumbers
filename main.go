package main

import (
	"bytes"
	"fmt"
	"github.com/disintegration/imaging"
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

	width := i[0].Int()
	height := i[1].Int()
	widthXScalar := i[2].Int()
	heightYScalar := i[3].Int()
	clusterCount := i[4].Int()
	quickMeans := i[5].Bool()
	kMeansTune := i[6].Float()
	srcArrayJS := i[7]
	outputBuffer := i[8]

	srcLen := srcArrayJS.Get("byteLength").Int()
	inputImageBytes := make([]uint8, srcLen)
	js.CopyBytesToGo(inputImageBytes, srcArrayJS)
	fmt.Printf("src bytes copied %f\n", time.Since(start).Seconds())
	fmt.Printf("Input Height: %d, width: %d, pixels: %d\n", height, width, len(inputImageBytes))

	//imgRgb, inputType, err := image.Decode(bytes.NewReader(inputImageBytes))
	imgRgb, err := imaging.Decode(bytes.NewReader(inputImageBytes), imaging.AutoOrientation(true))
	fmt.Printf("decode image done %f\n", time.Since(start).Seconds())
	if err != nil {
		panic("unable to read image")
	}
	fmt.Printf("Decode Height: %d, width: %d, pixels: %d\n", imgRgb.Bounds().Size().Y, imgRgb.Bounds().Size().X, imgRgb.Bounds().Size().X*imgRgb.Bounds().Size().Y)

	resizedImgRgb := resize.Resize(uint(imgRgb.Bounds().Size().X/widthXScalar), uint(imgRgb.Bounds().Size().Y/heightYScalar), imgRgb, resize.NearestNeighbor)
	fmt.Printf("image resize down %f\n", time.Since(start).Seconds())

	var colorPalette []clf.Color
	if quickMeans {
		fmt.Println("Running quick kmeans")
		colorPalette = pbn.DominantColors(resizedImgRgb, clusterCount, kMeansTune, false)
	} else {
		fmt.Println("Running full kmeans")
		colorPalette = pbn.DominantColorsAlt(imgRgb, clusterCount, int(kMeansTune))
	}
	fmt.Printf("color palette found, %d clusters, %f threshold %f\n", clusterCount, kMeansTune, time.Since(start).Seconds())

	snapImg := pbn.SnapColors(resizedImgRgb, colorPalette)
	fmt.Printf("snap done %f\n", time.Since(start).Seconds())

	newImage := resize.Resize(uint(imgRgb.Bounds().Size().X), uint(imgRgb.Bounds().Size().Y), snapImg, resize.NearestNeighbor)
	fmt.Printf("resize up done %f, widthX: %d, heightY: %d\n", time.Since(start).Seconds())

	imgBuf := new(bytes.Buffer)
	_ = png.Encode(imgBuf, newImage)
	bufBytes := imgBuf.Bytes()

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

	width := i[0].Int()
	height := i[1].Int()
	widthXScalar := i[2].Int()
	heightYScalar := i[3].Int()
	clusterCount := i[4].Int()
	quickMeans := i[5].Bool()
	kMeansTune := i[6].Float()
	srcArrayJS := i[7]

	srcLen := srcArrayJS.Get("byteLength").Int()
	inputImageBytes := make([]uint8, srcLen)
	js.CopyBytesToGo(inputImageBytes, srcArrayJS)
	fmt.Printf("src bytes copied %f\n", time.Since(start).Seconds())

	//imgRgb, inputType, err := image.Decode(bytes.NewReader(inputImageBytes))
	imgRgb, err := imaging.Decode(bytes.NewReader(inputImageBytes), imaging.AutoOrientation(true))
	fmt.Printf("decode image done %f\n", time.Since(start).Seconds())
	if err != nil {
		panic("unable to read image")
	}
	fmt.Printf("Height: %d, width: %d, pixels: %d\n", height, width, len(inputImageBytes))

	var colorPalette []clf.Color
	if quickMeans {
		resizedImgRgb := resize.Resize(uint(imgRgb.Bounds().Size().X/widthXScalar), uint(imgRgb.Bounds().Size().Y/heightYScalar), imgRgb, resize.NearestNeighbor)
		fmt.Printf("image resize down %f\n", time.Since(start).Seconds())

		fmt.Println("Running quick kmeans")
		colorPalette = pbn.DominantColors(resizedImgRgb, clusterCount, kMeansTune, false)
	} else {
		fmt.Println("Running full kmeans")
		colorPalette = pbn.DominantColorsAlt(imgRgb, clusterCount, int(kMeansTune))
	}
	fmt.Printf("color palette found, %d clusters, %f threshold %f\n", clusterCount, kMeansTune, time.Since(start).Seconds())

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

	//imgRgb, inputType, err := image.Decode(bytes.NewReader(inputImageBytes))
	imgRgb, err := imaging.Decode(bytes.NewReader(inputImageBytes), imaging.AutoOrientation(true))
	fmt.Printf("decode image done %f\n", time.Since(start).Seconds())
	if err != nil {
		panic("unable to read image")
	}
	fmt.Printf("Height: %d, width: %d, pixels: %d\n", height, width, len(inputImageBytes))

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

func stereogram(this js.Value, i []js.Value) interface{} {
	start := time.Now()

	heightYFactor := i[0].Int()
	widthXFactor := i[1].Int()
	patternWidthXFactor := i[2].Int()
	shiftAmplitude := i[3].Float()
	invert := i[4].Bool()
	srcArrayJS := i[5]
	outputBuffer := i[6]
	fmt.Printf("Args %s\n", i[:5])

	srcLen := srcArrayJS.Get("byteLength").Int()
	inputImageBytes := make([]uint8, srcLen)
	js.CopyBytesToGo(inputImageBytes, srcArrayJS)

	//baseImg := pbn.LoadImage("circle.png")
	//depthMapImgOrig := pbn.LoadImage("test_images/buns_dm.png")
	//baseImgColorOrig := pbn.LoadImage("test_images/buns.png")
	//depthMapImgOrig := pbn.LoadImage("color_bottle.png")
	depthMapImgOrig, err := imaging.Decode(bytes.NewReader(inputImageBytes), imaging.AutoOrientation(true))
	if err != nil {
		panic("unable to read image")
	}

	factoredWidthX := uint(depthMapImgOrig.Bounds().Dx() / widthXFactor)
	factoredWidthY := uint(depthMapImgOrig.Bounds().Dy() / heightYFactor)
	baseImg := resize.Resize(factoredWidthX, factoredWidthY, depthMapImgOrig, resize.NearestNeighbor)
	fmt.Printf("image resize down to %dx%d at %fs\n", factoredWidthX, factoredWidthY, time.Since(start).Seconds())

	//baseImg := pbn.LoadImage("test_images/gundam2.png")
	//depthMap := pbn.SimpleDepthMap(baseImg)
	//depthMap := pbn.GreyscaleDepthMap(baseImg)
	depthMap, _ := pbn.ColorDepthMap(baseImg)
	fmt.Printf("Depth Map array created at %fs\n", time.Since(start).Seconds())

	pattern := pbn.SimplePatternImage(baseImg.Bounds().Dx()/patternWidthXFactor, baseImg.Bounds().Dy())
	//pattern := pbn.PalettePatternImage(baseImg.Bounds().Dx()/10, baseImg.Bounds().Dy(), baseImgColorOrig)
	//pbn.SaveImage("pattern.png", pattern)
	fmt.Printf("Pattern created at %s %fs\n", pattern.Bounds(), time.Since(start).Seconds())

	stereogramImg := pbn.GenerateStereogram(depthMap, baseImg.Bounds().Dx(), baseImg.Bounds().Dy(), pattern, shiftAmplitude, invert)
	fmt.Printf("Stereogram created, size %s, at %fs\n", stereogramImg.Bounds(), time.Since(start).Seconds())

	imgBuf := new(bytes.Buffer)
	_ = png.Encode(imgBuf, stereogramImg)
	bufBytes := imgBuf.Bytes()

	js.CopyBytesToJS(outputBuffer, bufBytes)
	fmt.Printf("image (%d bytes) copied to outputBuffer %f\n", len(bufBytes), time.Since(start).Seconds())

	return js.ValueOf(len(imgBuf.Bytes()))
}

func makeDepthMap(this js.Value, i []js.Value) interface{} {
	start := time.Now()

	depthMapAlgo := i[0].Int()
	srcArrayJS := i[1]
	outputBuffer := i[2]

	println(depthMapAlgo)

	srcLen := srcArrayJS.Get("byteLength").Int()
	inputImageBytes := make([]uint8, srcLen)
	js.CopyBytesToGo(inputImageBytes, srcArrayJS)

	depthMapImgOrig, err := imaging.Decode(bytes.NewReader(inputImageBytes), imaging.AutoOrientation(true))
	if err != nil {
		panic("unable to read image")
	}
	fmt.Println(depthMapImgOrig.Bounds())

	var depthMapImage image.Image
	if depthMapAlgo == 1 {
		_, depthMapImage = pbn.GreyscaleDepthMap(depthMapImgOrig)
	} else if depthMapAlgo == 2 {
		_, depthMapImage = pbn.ColorDepthMap(depthMapImgOrig)
	} else {
		_, depthMapImage = pbn.SimpleDepthMap(depthMapImgOrig)
	}

	fmt.Printf("Depth Map created at %fs\n", time.Since(start).Seconds())

	imgBuf := new(bytes.Buffer)
	_ = png.Encode(imgBuf, depthMapImage)
	bufBytes := imgBuf.Bytes()

	js.CopyBytesToJS(outputBuffer, bufBytes)
	fmt.Printf("image (%d bytes) copied to outputBuffer %f\n", len(bufBytes), time.Since(start).Seconds())

	return js.ValueOf(len(imgBuf.Bytes()))
}

func registerCallbacks() {
	js.Global().Set("pixelizor", js.FuncOf(pixelizor))
	js.Global().Set("dominantColors", js.FuncOf(dominantColors))
	js.Global().Set("pixelizeFromPalette", js.FuncOf(pixelizeFromPalette))
	js.Global().Set("stereogram", js.FuncOf(stereogram))
	js.Global().Set("makeDepthMap", js.FuncOf(makeDepthMap))
}

func main() {
	c := make(chan struct{}, 0)

	println("WASM Go Initialized")
	// register functions
	registerCallbacks()
	<-c
}
