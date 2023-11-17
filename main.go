package main

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
)

func main() {
	var path string
	fmt.Print("Input file path: ")
	_, err := fmt.Scan(&path)
	if err != nil {
		fmt.Println("Error requesting input from user.")
		return
	}
	matrix, err := readImageToMatrix(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	//fmt.Println(matrix)
	_ = writeImageFromMatrix(matrix, path+"new")
}

// assertSignature asserts that a given os.File's signature (first n bytes of the file) are a given signature
// this is preferred to using mimetype as Go's built-in mimetype detection is quite deficient
func assertSignature(file *os.File, signature []byte) (bool, error) {
	fileSignature := make([]byte, len(signature)) // Make a new slice to hold the signature
	_, err := file.ReadAt(fileSignature, 0)       // read the first N bytes of the file
	if err != nil {
		return false, err
	}
	return bytes.Equal(fileSignature, signature), nil // Using bytes.Equal assert that the provided signature corresponds to the known signature
}

// readImageToMatrix takes a path for a given PNG image and returns a 2D uint8 slice that represents the matrix for that image
func readImageToMatrix(path string) ([][][4]uint32, error) {
	file, err := os.Open(path) // Open the provided file
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		_ = file.Close() // Ignore any errors resulting from closure
	}(file)

	isPNG, err := assertSignature(file, []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}) // Check the file's signature against the known PNG signature
	if err != nil {
		return nil, err
	}
	if !isPNG {
		return nil, errors.New("file is not a PNG image")
	}

	imageData, err := png.Decode(file) // Decode the os.File to an image.Image
	if err != nil {
		return nil, err
	}

	width, height := imageData.Bounds().Dx(), imageData.Bounds().Dy() // Get the width and height of the image

	// This operation is inherently 0(n²), alternative implementations have been studied where imageData was cast to an
	// image.NRGBA however this implementation didn't provide us with an operable 2D matrix

	matrix := Make2D[[4]uint32](width, height)
	fmt.Println(width, height)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, a := imageData.At(x, y).RGBA()
			// @TODO(Mauro): comment or explain why we use 257 instead of any other number
			matrix[x][y] = [4]uint32{r / 257, g / 257, b / 257, a / 257}
		}
	}

	return matrix, nil
}

// writeImageFromMatrix takes a matrix from readImageToMatrix and outputs a PNG image to a given path
func writeImageFromMatrix(matrix [][][4]uint32, path string) error {
	height := len(matrix[0])
	width := len(matrix)
	fmt.Println(width, height)
	img := image.NewNRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, a := matrix[x][y][0], matrix[x][y][1], matrix[x][y][2], matrix[x][y][3]
			img.SetNRGBA(x, y, color.NRGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)})
		}
	}

	// Create the file
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	// Encode as PNG
	return png.Encode(file, img)
}
