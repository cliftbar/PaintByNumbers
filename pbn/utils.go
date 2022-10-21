package pbn

import (
	"github.com/disintegration/imaging"
	"image"
	"image/png"
	"os"
	"path/filepath"
)

func LoadImage(imgPath string) image.Image {
	//existingImageFile, err := os.Open(filepath.Clean(imgPath))
	//if err != nil {
	//	panic(err)
	//}
	//defer existingImageFile.Close()

	//imageData, _, err := image.Decode(existingImageFile)
	imageData, err := imaging.Open(filepath.Clean(imgPath), imaging.AutoOrientation(true))
	if err != nil {
		panic(err)
	}

	return imageData
}

func SaveImage(path string, img image.Image) {
	outputFile, err := os.Create(path)
	if err != nil {
		// Handle error
	}

	err = png.Encode(outputFile, img)
	if err != nil {
		// Handle error
	}
}

// TODO check if I got width and height right
// ROW major
// mXn == rowXcol
func make2DSliceFloat64(yHeightRows int, xWidthCols int) [][]float64 {
	matrix := make([][]float64, xWidthCols)
	rows := make([]float64, yHeightRows*xWidthCols)
	for i := 0; i < xWidthCols; i++ {
		matrix[i] = rows[i*yHeightRows : (i+1)*yHeightRows]
	}

	return matrix
}
