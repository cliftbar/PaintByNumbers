package main

import (
	"bytes"
	"fmt"
	"github.com/muesli/clusters"
	"github.com/nfnt/resize"
	"image"
	_ "image/jpeg"
	"image/png"
	"paintByNumbers/pbn"
	"reflect"
	"syscall/js"
	"time"
	"unsafe"
)

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

func dominantColor(this js.Value, i []js.Value) interface{} {
	start := time.Now()
	//if len(i)+1 < 3 {
	//	fmt.Println("Not enough arguments")
	//} else {
	//	fmt.Println("Calculating colors")
	//}
	fmt.Println("Calculating colors")

	height := i[0].Int()
	width := i[1].Int()
	widthXScalar := 5
	heightYScalar := 5
	clusterCount := 5

	array := i[2]
	fmt.Println(array.Get("byteLength"))
	pixelCount := height * width
	srcLen := array.Get("byteLength").Int()
	inputImageBytes := make([]uint8, srcLen)
	js.CopyBytesToGo(inputImageBytes, array)

	imgRgb, inputType, err := image.Decode(bytes.NewReader(inputImageBytes))
	if err != nil {
		panic("unable to read image")
	}

	fmt.Printf("byteLen: %d, Height: %d, width: %d, pixels: %d, type: %s\n", srcLen, height, width, pixelCount, inputType)

	resizedImgRgb := resize.Resize(uint(imgRgb.Bounds().Size().X/widthXScalar), uint(imgRgb.Bounds().Size().Y/heightYScalar), imgRgb, resize.NearestNeighbor)
	//resizedImgRgb := <-asyncResize(imgRgb, uint(imgRgb.Bounds().Size().X/widthXScalar), uint(imgRgb.Bounds().Size().Y/heightYScalar))
	fmt.Println("image loaded")
	colorPalette := pbn.DominantColors(resizedImgRgb, clusterCount, 0.01, false)
	//colorPalette := <-asyncDominantColors(resizedImgRgb, clusterCount, 0.01)
	fmt.Println("color palette found")
	snapImg := pbn.SnapColors(resizedImgRgb, colorPalette)
	//snapImg := <-asyncSnapColors(resizedImgRgb, colorPalette)

	fmt.Println("snap done")

	//newImage := resize.Resize(uint(imgRgb.Bounds().Size().X), uint(imgRgb.Bounds().Size().Y), snapImg, resize.NearestNeighbor)
	newImage := <-asyncResize(snapImg, uint(imgRgb.Bounds().Size().X), uint(imgRgb.Bounds().Size().Y))

	imgBuf := new(bytes.Buffer)
	_ = png.Encode(imgBuf, newImage)
	bufBytes := imgBuf.Bytes()

	buffHeader := (*reflect.SliceHeader)(unsafe.Pointer(&bufBytes))

	retMap := map[string]interface{}{
		"imgPtr":    buffHeader.Data,
		"imgPtrLen": buffHeader.Len,
	}

	fmt.Printf("Wasm Done %fs\n", time.Since(start).Seconds())
	return retMap
}

//func updateImage(img *image.RGBA, start time.Time) {
//	enc := imgio.JPEGEncoder(90)
//	err := enc(&s.outBuf, img)
//	if err != nil {
//		s.log(err.Error())
//		return
//	}
//
//	dst := js.Global().Get("Uint8Array").New(len(s.outBuf.Bytes()))
//	n := js.CopyBytesToJS(dst, s.outBuf.Bytes())
//	s.console.Call("log", "bytes copied:", strconv.Itoa(n))
//	js.Global().Call("displayImage", dst)
//	s.console.Call("log", "time taken:", time.Now().Sub(start).String())
//	s.outBuf.Reset()
//}

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
