package manipulations

import (
	"errors"
	"matrix-image-manipulation/utils"
)

// ConvertToGreyScale converts a matrix to greyscale equivalent
func ConvertToGreyScale(matrix [][][4]uint32) ([][][4]uint32, error) {
	height := len(matrix) // Get the height of the matrix
	if height == 0 {
		return nil, errors.New("empty matrix")
	}

	width := len(matrix[0])

	greyScaleImage := utils.Make2D[[4]uint32](height, width) // Create a new image

	// Iterate through all pixels
	for y := 0; y < height; y++ {
		for x := range matrix[y] {
			current := matrix[y][x]                                                     // Convenience, cleans up the following line, does allocate more memory however
			r, g, b, a := current[0], current[1], current[2], current[3]                // Thank you Go for not providing list expansion
			luminance := uint32(float64(r)*0.299 + float64(g)*0.587 + float64(b)*0.114) // Apply  the formula
			greyScaleImage[y][x] = [4]uint32{luminance, luminance, luminance, a}        // Write the luminance value, keeping the alpha as current
		}
	}
	return greyScaleImage, nil
}
