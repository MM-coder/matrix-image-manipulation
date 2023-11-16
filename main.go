package main

import (
	"bytes"
	"errors"
	"fmt"
	"image/png"
	"os"
)

func main() {
	matrix, err := readImageToMatrix("./shrekfly.png")
	if err != nil {
		return
	}
	fmt.Println(matrix)
	fmt.Println(err)
}

func assertSignature(file *os.File, signature []byte) (bool, error) {
	fileSignature := make([]byte, len(signature))
	_, err := file.ReadAt(fileSignature, 0)
	if err != nil {
		return false, err
	}
	return bytes.Equal(fileSignature, signature), nil
}

func readImageToMatrix(path string) ([][]uint8, error) {
	file, err := os.Open(path) // Open the provided file
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		_ = file.Close() // Ignore any errors resulting from closure
	}(file)

	isPNG, err := assertSignature(file, []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A})
	if err != nil {
		return nil, err
	}
	if !isPNG {
		return nil, errors.New("file is not a PNG image")
	}

	imageData, err := png.Decode(file)
	if err != nil {
		return nil, err
	}

	width, height := imageData.Bounds().Dx(), imageData.Bounds().Dy()

	// This operation is inherently 0(n**2), alternative implementations have been studied where imageData was cast to a
	// image.NRGBA however this implementation didn't provide us with an operable 2D matrix

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			colour := imageData.At(x, y)
			fmt.Println(colour)
		}
	}

	return nil, errors.New("unable to read image file as NRGBA")

}
