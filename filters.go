package matrix_image_manipulation

import (
	"errors"
	"math"
)

// gaussianFilter applies a Gaussian filter to a matrix representing an image
func gaussianFilter(matrix [][][4]uint32, kernelSize int, sigma float64) ([][][4]uint32, error) {
	height := len(matrix) // Get the height of the matrix
	if height == 0 {
		return nil, errors.New("empty matrix")
	}
	width := len(matrix[0])

	kernel := generateGaussianKernel(kernelSize, sigma)

	filteredMatrix := Make2D[[4]uint32](height, width)

	// Apply the Gaussian kernel to each pixel
	kOffset := kernelSize / 2
	for y := kOffset; y < height-kOffset; y++ {
		for x := kOffset; x < width-kOffset; x++ {
			filteredMatrix[y][x] = applyKernel(x, y, matrix, kernel, kOffset)
		}
	}

	return filteredMatrix, nil
}

// applyKernel applies the given Guassian kernel to a single pixel
func applyKernel(x int, y int, matrix [][][4]uint32, kernel [][]float64, kOffset int) [4]uint32 {
	var sum [4]float64
	for ky := 0; ky < len(kernel); ky++ {
		for kx := 0; kx < len(kernel); kx++ {
			px := matrix[y+kOffset-ky][x+kOffset-kx]
			for i := range px {
				sum[i] += float64(px[i]) * kernel[ky][kx]
			}
		}
	}

	var result [4]uint32
	for i := range result {
		result[i] = uint32(math.Min(math.Max(sum[i], 0), 255))
	}
	return result
}

// generateGaussianKernel generates a Gaussian kernel
func generateGaussianKernel(size int, sigma float64) [][]float64 {
	kernel := Make2D[float64](size, size)

	sum := 0.0
	offset := size / 2

	for y := -offset; y <= offset; y++ {
		for x := -offset; x <= offset; x++ {
			val := (1.0 / (2.0 * math.Pi * sigma * sigma)) * math.Exp(-(float64(x*x+y*y) / (2.0 * sigma * sigma)))
			kernel[y+offset][x+offset] = val
			sum += val
		}
	}

	// Normalize the kernel
	for y := range kernel {
		for x := range kernel[y] {
			kernel[y][x] /= sum
		}
	}

	return kernel
}
