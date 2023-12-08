package manipulations

import (
	"errors"
	"math"
	"matrix-image-manipulation/utils"
)

// GaussianFilter applies a Gaussian filter to a matrix representing an image.
func GaussianFilter(matrix [][][4]uint32, kernelSize int, sigma float64) ([][][4]uint32, error) {
	height := len(matrix) // Get the height of the matrix
	if height == 0 {
		return nil, errors.New("empty matrix")
	}
	width := len(matrix[0])

	// Generate the Gaussian kernel with the given size and standard deviation (sigma).
	kernel := generateGaussianKernel(kernelSize, sigma)

	filteredMatrix := utils.Make2D[[4]uint32](height, width)

	// Apply the Gaussian kernel to each pixel.
	// kOffset is used to handle border effects by avoiding out-of-bounds indices.
	kOffset := kernelSize / 2
	for y := kOffset; y < height-kOffset; y++ {
		for x := kOffset; x < width-kOffset; x++ {
			// Apply the kernel to the pixel at (x, y) and store the result.
			filteredMatrix[y][x] = applyKernel(x, y, matrix, kernel, kOffset)
		}
	}

	return filteredMatrix, nil
}

// applyKernel applies the given Gaussian kernel to a single pixel.
func applyKernel(x int, y int, matrix [][][4]uint32, kernel [][]float64, kOffset int) [4]uint32 {
	var sum [4]float64
	for ky := 0; ky < len(kernel); ky++ {
		for kx := 0; kx < len(kernel); kx++ {
			// Multiply each kernel coefficient with the corresponding pixel value.
			px := matrix[y+kOffset-ky][x+kOffset-kx]
			for i := range px {
				sum[i] += float64(px[i]) * kernel[ky][kx]
			}
		}
	}

	// Convert the summed values back to uint32, ensuring they remain within the valid range [0, 255].
	var result [4]uint32
	for i := range result {
		result[i] = uint32(math.Min(math.Max(sum[i], 0), 255))
	}
	return result
}

// generateGaussianKernel generates a Gaussian kernel for image blurring.
// The Gaussian kernel is a square matrix used for the blurring effect.
func generateGaussianKernel(size int, sigma float64) [][]float64 {
	kernel := utils.Make2D[float64](size, size)

	sum := 0.0
	offset := size / 2

	// Fill the kernel with values computed using the Gaussian function.
	// The kernel values are based on the distance from the center, modulated by the sigma (standard deviation).
	for y := -offset; y <= offset; y++ {
		for x := -offset; x <= offset; x++ {
			val := (1.0 / (2.0 * math.Pi * sigma * sigma)) * math.Exp(-(float64(x*x+y*y) / (2.0 * sigma * sigma)))
			kernel[y+offset][x+offset] = val
			sum += val
		}
	}

	// Normalize the kernel so that the sum of all its values equals 1.
	// This ensures that applying the kernel to an image preserves the image's brightness.
	for y := range kernel {
		for x := range kernel[y] {
			kernel[y][x] /= sum
		}
	}

	return kernel
}
