package utils

import (
	"errors"
	"image"
	"image/color"
	"image/png"
	"os"
)

// ReadImageToMatrix takes a path for a given PNG image and returns a 2D uint8 slice that represents the matrix for that image
func ReadImageToMatrix(path string) ([][][4]uint32, error) {
	file, err := os.Open(path) // Open the provided file
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		_ = file.Close() // Ignore any errors resulting from closure
	}(file)

	isPNG, err := AssertSignature(file, []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}) // Check the file's signature against the known PNG signature
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

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, a := imageData.At(x, y).RGBA()
			// @NOTE(Mauro):In Go, the image.Image interface's At(x, y).RGBA() method returns color values as four
			// uint32s, where each color component (Red, Green, Blue, Alpha) is represented in the range [0, 65535].
			// This is essentially a 16-bit per channel representation.
			//
			// However, typical image processing and display systems use 8-bit per channel representations,
			// where each color component is in the range [0, 255].
			//
			// 257 is used (which is 65535/255) as the divisor is to properly scale down these 16-bit values
			// to 8-bit values.
			// When you divide each component by 257, you map the 16-bit range [0, 65535] to the 8-bit range [0, 255]
			// effectively.
			matrix[x][y] = [4]uint32{r / 257, g / 257, b / 257, a / 257}
		}
	}

	return matrix, nil
}

// WriteImageFromMatrix takes a matrix from ReadImageToMatrix and outputs a PNG image to a given path
func WriteImageFromMatrix(matrix [][][4]uint32, path string) error {
	height, width := len(matrix[0]), len(matrix)           // Get the width and height of the matrix
	img := image.NewNRGBA(image.Rect(0, 0, width, height)) // Create a new generic image with the appropriate width and height

	// Iterate through the matrix and set values for each of the pixels
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			current := matrix[x][y]
			r, g, b, a := current[0], current[1], current[2], current[3]
			img.SetNRGBA(x, y, color.NRGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)})
		}
	}

	file, err := os.Create(path) // Create the file
	if err != nil {              // Return errors
		return err
	}
	defer func(file *os.File) {
		_ = file.Close() // Close the file
	}(file)

	// Encode as PNG
	return png.Encode(file, img)
}
