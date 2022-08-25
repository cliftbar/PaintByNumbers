package main

// tinygo main
import (
	"bytes"
	"fmt"
	clf "github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/clusters"
	"github.com/nfnt/resize"
	_ "golang.org/x/image/webp"
	"image"
	_ "image/jpeg"
	"image/png"
	"paintByNumbers/pbn"
	"syscall/js"
	"time"
)

// var imgBuf *bytes.Buffer
var inputImageBytes []uint8

// Pin buffer to global, so it doesn't get GC'd
func asyncResize(imgRgb image.Image, widthX uint, heightY uint) <-chan image.Image {
	r := make(chan image.Image)

	go func() {
		defer close(r)
		r <- resize.Resize(widthX, heightY, imgRgb, resize.NearestNeighbor)
	}()

	return r
}

func asyncDominantColors(imgRgb image.Image, clusterCount int, deltaThreshold float64) <-chan clusters.Clusters {
	r := make(chan clusters.Clusters)

	go func() {
		defer close(r)
		fmt.Println("start dom colors")
		r <- pbn.DominantColors(imgRgb, clusterCount, deltaThreshold, false)
	}()

	return r
}
func asyncSnapColors(imgRgb image.Image, colorPalette clusters.Clusters) <-chan image.Image {
	r := make(chan image.Image)

	go func() {
		defer close(r)
		r <- pbn.SnapColors(imgRgb, colorPalette)
	}()

	return r
}

// func dominantColor(height int, width int, inputImageBytes []byte) interface{} {
func dominantColor(this js.Value, i []js.Value) interface{} {
	start := time.Now()
	//if len(i)+1 < 3 {
	//	fmt.Println("Not enough arguments")
	//} else {
	//	fmt.Println("Calculating colors")
	//}

	//fmt.Println(i[0])
	//fmt.Println(i[1])
	//fmt.Println(i[2])
	height := i[0].Int()
	width := i[1].Int()
	widthXScalar := i[2].Int()
	heightYScalar := i[3].Int()
	clusterCount := i[4].Int()
	//
	array := i[5]
	//fmt.Println(array.Get("byteLength"))
	//pixelCount := height * width
	srcLen := array.Get("byteLength").Int()
	inputImageBytes = make([]uint8, srcLen)
	js.CopyBytesToGo(inputImageBytes, array)
	//inputImageBytes := make([]byte, 10)
	//fmt.Println(inputImageBytes)
	fmt.Printf("decode image %f\n", time.Since(start).Seconds())
	imgRgb, inputType, err := image.Decode(bytes.NewReader(inputImageBytes))
	fmt.Printf("decode image done %f\n", time.Since(start).Seconds())
	if err != nil {
		panic("unable to read image")
	}

	fmt.Printf("Height: %d, width: %d, pixels: %d, type: %s\n", height, width, len(inputImageBytes), inputType)
	fmt.Printf("image resize %f\n", time.Since(start).Seconds())
	resizedImgRgb := resize.Resize(uint(imgRgb.Bounds().Size().X/widthXScalar), uint(imgRgb.Bounds().Size().Y/heightYScalar), imgRgb, resize.NearestNeighbor)
	//resizedImgRgb := <-asyncResize(imgRgb, uint(imgRgb.Bounds().Size().X/widthXScalar), uint(imgRgb.Bounds().Size().Y/heightYScalar))
	fmt.Printf("Calculating colors %f\n", time.Since(start).Seconds())
	colorPalette := pbn.DominantColors(resizedImgRgb, clusterCount, 0.01, false)
	//colorPalette := <-asyncDominantColors(resizedImgRgb, clusterCount, 0.01)
	fmt.Printf("color palette found %f\n", time.Since(start).Seconds())
	snapImg := pbn.SnapColors(resizedImgRgb, colorPalette)
	//snapImg := <-asyncSnapColors(resizedImgRgb, colorPalette)

	fmt.Printf("snap done %f\n", time.Since(start).Seconds())

	newImage := resize.Resize(uint(imgRgb.Bounds().Size().X), uint(imgRgb.Bounds().Size().Y), snapImg, resize.NearestNeighbor)
	fmt.Printf("resize up done %f\n", time.Since(start).Seconds())
	//newImage := <-asyncResize(snapImg, uint(imgRgb.Bounds().Size().X), uint(imgRgb.Bounds().Size().Y))
	var imgBuf *bytes.Buffer
	var bufBytes []uint8

	imgBuf = new(bytes.Buffer)
	_ = png.Encode(imgBuf, newImage)
	bufBytes = imgBuf.Bytes()

	//buffHeader := unsafe.Pointer(&bufBytes)
	output := i[6]
	js.CopyBytesToJS(output, bufBytes)
	fmt.Printf("output copy done %f\n", time.Since(start).Seconds())

	colorPaletteHex := make([]string, 0)
	colorPaletteHexStr := ""
	for _, centroid := range colorPalette {
		hex := clf.Color{
			R: centroid.Center[0] / 255.0,
			G: centroid.Center[1] / 255.0,
			B: centroid.Center[2] / 255.0,
		}.Hex()
		colorPaletteHex = append(colorPaletteHex, hex)
		if colorPaletteHexStr == "" {
			colorPaletteHexStr = hex
		} else {
			colorPaletteHexStr = colorPaletteHexStr + "," + hex
		}
	}

	//retMap := map[string]interface{}{
	//	//"imgPtr":    js.ValueOf(bufBytes),
	//	"colors": colorPaletteHexStr,
	//}

	//js.CopyBytesToJS(, bufBytes)

	fmt.Printf("Wasm Done %fs\n", time.Since(start).Seconds())
	return js.ValueOf(colorPaletteHexStr)
}

func registerCallbacks() {
	js.Global().Set("dominantColor", js.FuncOf(dominantColor))
}

func main() {
	c := make(chan struct{}, 0)

	println("WASM Go Initialized")
	// register functions
	registerCallbacks()
	<-c
}
