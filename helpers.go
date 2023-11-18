package matrix_image_manipulation

import (
	"bytes"
	"os"
)

// Make2D makes a 2D slice of any type of the given width and height
// src:  https://stackoverflow.com/a/71781206 (adapted)
func Make2D[Type any](n, m int) [][]Type {
	matrix := make([][]Type, n)
	rows := make([]Type, n*m)
	for i, startRow := 0, 0; i < n; i, startRow = i+1, startRow+m {
		endRow := startRow + m
		matrix[i] = rows[startRow:endRow:endRow]
	}
	return matrix
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
