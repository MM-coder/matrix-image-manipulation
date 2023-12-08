package manipulations

import (
	"errors"
	"math"
	"matrix-image-manipulation/utils"
)

// AdjustContrast alters the contrast of an image using the formula g(u) = mu*u + b.
func AdjustContrast(matrix [][][4]uint32, m, b float64) ([][][4]uint32, error) {
	height := len(matrix) // Get the height of the matrix
	if height == 0 {
		return nil, errors.New("empty matrix")
	}

	width := len(matrix[0])

	// Create a new matrix for the contrast-adjusted image.
	contrastMatrix := utils.Make2D[[4]uint32](height, width)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			var adjustedPixel [4]uint32
			originalPixel := matrix[y][x]

			for i := 0; i < 3; i++ { // Iterate over R, G, B components (not A)
				// Apply the contrast formula.
				// Clamp the result to the range [0, 255].
				adjustedValue := math.Max(math.Min(m*float64(originalPixel[i])+b, 255), 0)
				adjustedPixel[i] = uint32(adjustedValue)
			}
			adjustedPixel[3] = originalPixel[3] // Preserve the alpha channel

			contrastMatrix[y][x] = adjustedPixel
		}
	}

	return contrastMatrix, nil
}

// AdjustLuminosity alters the contrast of an image using the formula g(u) = mu*u + b.
func AdjustLuminosity(matrix [][][4]uint32, b float64) ([][][4]uint32, error) {
	height := len(matrix) // Get the height of the matrix
	if height == 0 {
		return nil, errors.New("empty matrix")
	}

	width := len(matrix[0])

	// Create a new matrix for the contrast-adjusted image.
	contrastMatrix := utils.Make2D[[4]uint32](height, width)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			var adjustedPixel [4]uint32
			originalPixel := matrix[y][x]

			for i := 0; i < 3; i++ { // Iterate over R, G, B components (not A)
				// Apply the luminosity formula.
				// Clamp the result to the range [0, 255].
				adjustedValue := math.Max(math.Min(float64(originalPixel[i])+b, 255), 0)
				adjustedPixel[i] = uint32(adjustedValue)
			}
			adjustedPixel[3] = originalPixel[3] // Preserve the alpha channel

			contrastMatrix[y][x] = adjustedPixel
		}
	}

	return contrastMatrix, nil
}
