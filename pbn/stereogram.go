package pbn

import (
	"github.com/disintegration/imaging"
	"image"
	"image/color"
	"math/rand"
)

func SimplePatternImage(widthX int, heightY int) image.Image {
	img := imaging.New(widthX, heightY, color.Black)

	//blue := color.RGBA{
	//	R: 0,
	//	G: 0,
	//	B: 255,
	//	A: 255,
	//}

	for x := 0; x < img.Bounds().Dx(); x++ {
		for y := 0; y < img.Bounds().Dy(); y++ {

			//if 10 < x && x < 15 && 10 < y && y < 15 {
			//	img.Set(x, y, )
			//} else {
			//	img.Set(x, y, color.Gray{Y: uint8(rand.Intn(255))})
			//}
			img.Set(x, y, color.Gray{Y: uint8(rand.Intn(255))})
		}
	}

	return img
}

func PalettePatternImage(widthX int, heightY int, src image.Image) image.Image {
	//palette := DominantColors(src, 5, 0.1, false)
	palette := DominantColorsAlt(src, 5, 10)
	img := imaging.New(widthX, heightY, color.Black)

	for x := 0; x < img.Bounds().Size().X; x++ {
		for y := 0; y < img.Bounds().Size().Y; y++ {

			if 10 < x && x < 15 && 10 < y && y < 15 {
				img.Set(x, y, color.RGBA{R: 0, G: 0, B: 255, A: 255})
			} else {
				img.Set(x, y, palette[rand.Intn(len(palette))])
			}
		}
	}

	return img
}

func SimpleDepthMap(src image.Image) [][]float64 {
	outMap := make2DSliceFloat64(src.Bounds().Dy(), src.Bounds().Dx())
	depthMapImg := imaging.New(src.Bounds().Dx(), src.Bounds().Dy(), color.Black)

	for x := 0; x < src.Bounds().Dx(); x++ {
		for y := 0; y < src.Bounds().Dy(); y++ {

			R, G, B, _ := src.At(x, y).RGBA()
			//r, g, b := colorful.LinearRgb(float64(R), float64(G), float64(B)).RGB255()
			pixGrayValue := ((float32(R) * 0.3) + (float32(G) * 0.59) + (float32(B) * .11)) / 65535 * 255
			//pixGrayValue := (float64(r) * 0.3) + (float64(g) * 0.59) + (float64(b) * .11)

			snapGrey := 0.0
			snapGreyPix := 0
			if pixGrayValue > 128 {
				snapGrey = 1.0
				snapGreyPix = 255
			}

			depthMapImg.Set(x, y, color.Gray{Y: uint8(snapGreyPix)})
			outMap[x][y] = snapGrey
		}
	}

	SaveImage("depthMap.png", depthMapImg)

	return outMap
}

func ColorDepthMap(src image.Image) [][]float64 {
	outMap := make2DSliceFloat64(src.Bounds().Dy(), src.Bounds().Dx())
	depthMapImg := imaging.New(src.Bounds().Dx(), src.Bounds().Dy(), color.Black)
	max := 0.0

	for x := 0; x < src.Bounds().Dx(); x++ {
		for y := 0; y < src.Bounds().Dy(); y++ {

			R, G, B, _ := src.At(x, y).RGBA()
			//r, g, b := colorful.LinearRgb(float64(R), float64(G), float64(B)).RGB255()
			pixGrayValue := ((float64(R) * 0.3) + (float64(G) * 0.59) + (float64(B) * .11)) / 65535.0 * 255.0
			//pixGrayValue := (float64(r) * 0.3) + (float64(g) * 0.59) + (float64(b) * .11)
			if max < pixGrayValue {
				max = pixGrayValue
			}

			depthMapImg.Set(x, y, color.Gray{Y: uint8(pixGrayValue)})
			outMap[x][y] = pixGrayValue
		}
	}

	for x := 0; x < src.Bounds().Dx(); x++ {
		for y := 0; y < src.Bounds().Dy(); y++ {
			outMap[x][y] = outMap[x][y] / max
		}
	}
	//println(max)

	//SaveImage("depthMap.png", depthMapImg)

	return outMap
}

func GreyscaleDepthMap(src image.Image) [][]float64 {
	outMap := make2DSliceFloat64(src.Bounds().Dy(), src.Bounds().Dx())
	depthMapImg := imaging.New(src.Bounds().Dx(), src.Bounds().Dy(), color.Black)
	max := 0.0

	for x := 0; x < src.Bounds().Dx(); x++ {
		for y := 0; y < src.Bounds().Dy(); y++ {

			R, _, _, _ := src.At(x, y).RGBA()
			//r, g, b := colorful.LinearRgb(float64(R), float64(G), float64(B)).RGB255()
			pixGrayValue := float64(R) / 65535.0
			//pixGrayValue := (float64(r) * 0.3) + (float64(g) * 0.59) + (float64(b) * .11)

			//snapGrey := 0.0
			//snapGreyPix := 0
			//if pixGrayValue > 128 {
			//	snapGrey = 1.0
			//	snapGreyPix = 255
			//}
			//if pixGrayValue != 0 {
			//	print(pixGrayValue)
			//}

			if max < pixGrayValue {
				max = pixGrayValue
			}

			depthMapImg.Set(x, y, color.Gray{Y: uint8(pixGrayValue * 255)})
			outMap[x][y] = pixGrayValue
		}
	}

	for x := 0; x < src.Bounds().Dx(); x++ {
		for y := 0; y < src.Bounds().Dy(); y++ {
			outMap[x][y] = outMap[x][y] / float64(max)
		}
	}

	println(max)

	SaveImage("depthMap.png", depthMapImg)

	return outMap
}

// amplitude 0.1
func GenerateStereogram(depthMap [][]float64, width int, height int, pattern image.Image, shiftAmplitude float64, invert bool) image.Image {
	autostereogram := imaging.New(width, height, color.Black)
	for yHeightRow := 0; yHeightRow < height; yHeightRow++ {
		for xWidthCol := 0; xWidthCol < width; xWidthCol++ {

			if xWidthCol < pattern.Bounds().Dx() {
				autostereogram.Set(xWidthCol, yHeightRow, pattern.At(xWidthCol, yHeightRow%pattern.Bounds().Dy()))
			} else {
				//shift = int(depthmap[r, c] * shift_amplitude * pattern.shape[1])
				//autostereogram[r, c] = autostereogram[r, c - pattern.shape[1] + shift]
				depth := depthMap[xWidthCol][yHeightRow]
				if invert {
					depth = 1 - depth
				}
				shift := int(depth * shiftAmplitude * float64(pattern.Bounds().Dx()))

				newPosX := xWidthCol + shift - pattern.Bounds().Dx()
				pixColor := autostereogram.At(newPosX, yHeightRow)

				autostereogram.Set(xWidthCol, yHeightRow, pixColor)
			}

			//shift := int(depthMap[xWidthCol][yHeightRow] * shiftAmplitude * float64(pattern.Bounds().Dx()))
			//pixColor := pattern.At(xWidthCol, yHeightRow)
			//newPosX := xWidthCol + shift - pattern.Bounds().Dx()
			//autostereogram.Set(newPosX, yHeightRow, pixColor)
			//autostereogram.Set(xWidthCol, yHeightRow, pixColor)
		}
	}

	return autostereogram
}
