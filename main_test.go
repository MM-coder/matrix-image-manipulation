package main

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"math/rand"
	"matrix-image-manipulation/manipulations"
	"matrix-image-manipulation/utils"
	"os"
	"reflect"
	"strconv"
	"testing"
)

// loadImage handles loading an image from a given filePath.
func loadImage(filePath string) (image.Image, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	img, _, err := image.Decode(file)
	return img, err
}

// assertColourEquality checks the equality of two color.Color objects
// we're required to do this as RGBA and NRGBA are different for Go, even though their representation is the same
func assertColourEquality(colour1 color.Color, colour2 color.Color) bool {
	r1, g1, b1, a1 := colour1.RGBA()
	r2, g2, b2, a2 := colour2.RGBA()
	return r1 == r2 && g1 == g2 && b1 == b2 && a1 == a2
}

// generateRandomImage generates a random valid image of a given width and height.
func generateRandomImage(width int, height int) [][][4]uint32 {
	img := utils.Make2D[[4]uint32](height, width)

	// Iterate through all elements of the list and assign it to a random array with 4 random values between 0-255
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img[y][x] = [4]uint32{
				uint32(rand.Intn(256)), // Red
				uint32(rand.Intn(256)), // Green
				uint32(rand.Intn(256)), // Blue
				uint32(rand.Intn(256)), // Alpha
			}
		}
	}
	return img
}

// assertValidMatrix asserts that a given "matrix" is valid, based only on the RGBA values being within the valid interval
func assertValidMatrix(matrix [][][4]uint32) error {
	// Iterate through all the elements of the matrix
	for y := range matrix {
		for x := range matrix[y] {
			pixel := matrix[y][x]
			if pixel[0] > 255 || pixel[1] > 255 || pixel[2] > 255 || pixel[3] > 255 { // If any of the values are over 255, raise an error
				return errors.New("invalid matrix: values exceed 255")
			}
		}
	}
	return nil
}

// TestMatrixContinuity tests that if an image is converted to a matrix, and then back, it remains the same
func TestMatrixContinuity(t *testing.T) {

	var testImagePath = ".github/test_images/gnome.png"

	matrix, err := utils.ReadImageToMatrix(testImagePath)
	if err != nil {
		t.Fatalf("Failed to load the test image: %s", err)
		return
	}

	temporaryPath := t.TempDir() + "gnome.png" // create a temporary path for the output image

	_ = utils.WriteImageFromMatrix(matrix, temporaryPath)

	generatedImage, err := loadImage(temporaryPath)
	if err != nil {
		t.Fatalf("Failed loading the created image: %s", err)
		return
	}

	originalImage, err := loadImage(testImagePath)
	if err != nil {
		t.Fatalf("Failed loading the created image: %s", err)
		return
	}

	if generatedImage.Bounds() != originalImage.Bounds() {
		t.Errorf("Image dimensions are different: %s != %s", generatedImage.Bounds(), originalImage.Bounds())
		return
	}

	// Iterate through all the image's pixels and assert they are the same
	for y := originalImage.Bounds().Min.Y; y < originalImage.Bounds().Max.Y; y++ {
		for x := originalImage.Bounds().Min.X; x < originalImage.Bounds().Max.X; x++ {
			if !assertColourEquality(originalImage.At(x, y), generatedImage.At(x, y)) {
				fmt.Println(originalImage.At(x, y), generatedImage.At(x, y))
				t.Errorf("Image pixels are not the same, failed at coord (%s, %s)", strconv.Itoa(x), strconv.Itoa(y))
				return
			}
		}
	}

}

// TestMake2D tests the Make2D function, asserting it can correctly make a nxm matrix of type [4]int32 which is our use-case
func TestMake2D(t *testing.T) {
	width, height := rand.Intn(100)+1, rand.Intn(100)+1 // As rand.Intn returns from [0, n] we must add one to assert n is never 0
	matrix := utils.Make2D[[4]int32](width, height)

	// Check if the number of rows is n
	if len(matrix) != width {
		t.Errorf("Expected %d rows, got %d", width, len(matrix))
		return
	}

	// Check if each row has m elements
	for i, row := range matrix {
		if len(row) != height {
			t.Errorf("Expected %d elements in row %d, got %d", height, i, len(row))
			return
		}
	}

	// Check if each element's capacity is 4, which is as requirement
	for _, row := range matrix {
		for _, elem := range row {
			if !(cap(elem) == 4) {
				t.Errorf("Expected element to be initialized to an empty slice, got %v", elem)
				return
			}
		}
	}
}

// TestConvertToGreyScaleWithRandomInput tests convertToGreyScale function with a randomly generated input.
func TestConvertToGreyScaleWithRandomInput(t *testing.T) {

	width, height := rand.Intn(3841), rand.Intn(2161) // Random dimensions up to 4k
	randomImage := generateRandomImage(width, height)

	expected := make([][][4]uint32, height)
	for y := 0; y < height; y++ {
		expected[y] = make([][4]uint32, width)
		for x := 0; x < width; x++ {
			r, g, b, a := randomImage[y][x][0], randomImage[y][x][1], randomImage[y][x][2], randomImage[y][x][3]
			luminance := uint32(float64(r)*0.299 + float64(g)*0.587 + float64(b)*0.114)
			expected[y][x] = [4]uint32{luminance, luminance, luminance, a}
		}
	}

	result, err := manipulations.ConvertToGreyScale(randomImage)
	if err != nil {
		t.Errorf("convertToGreyScale() returned an unexpected error: %v", err)
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("convertToGreyScale() result does not match expected output")
	}
}

// TestGenerateRandomImage tests the generateRandomImage function.
func TestGenerateRandomImage(t *testing.T) {
	width, height := rand.Intn(3841), rand.Intn(2161) // Random dimensions up to 4k

	// Generate a random image
	randomImage := generateRandomImage(width, height)

	// Assert that the generated image is valid
	err := assertValidMatrix(randomImage)
	if err != nil {
		t.Errorf("generateRandomImage() generated an invalid matrix: %v", err)
	}
}

// TestAssertValidMatrix tests the assertValidMatrix function with a random number of valid and invalid matrices.
func TestAssertValidMatrix(t *testing.T) {

	// Inline function to generate a random pixel
	generateRandomPixel := func(maxVal uint32) [4]uint32 {
		return [4]uint32{
			uint32(rand.Intn(int(maxVal) + 1)),
			uint32(rand.Intn(int(maxVal) + 1)),
			uint32(rand.Intn(int(maxVal) + 1)),
			uint32(rand.Intn(int(maxVal) + 1)),
		} // @NOTE(Mauro): I miss list comprehensions
	}

	// Function to generate a random matrix using Make2D
	generateRandomMatrix := func(width int, height int, valid bool) [][][4]uint32 {
		matrix := utils.Make2D[[4]uint32](height, width)
		maxVal := uint32(255)
		if !valid {
			maxVal = 300 // Ensure an invalid matrix
		}
		for i := range matrix {
			for j := range matrix[i] {
				matrix[i][j] = generateRandomPixel(maxVal)
			}
		}
		return matrix
	}

	// Generate a random number of test cases
	numTests := rand.Intn(10) + 1 // At least 1 test, up to 10

	for i := 0; i < numTests; i++ {

		valid := rand.Intn(2) == 0                        // Randomly decide if this matrix should be valid or not
		width, height := rand.Intn(3841), rand.Intn(2161) // Random dimensions up to 4k
		matrix := generateRandomMatrix(width, height, valid)

		err := assertValidMatrix(matrix)
		if valid && err != nil {
			t.Errorf("assertValidMatrix() returned an error for a valid matrix: %v", err)
			return
		} else if !valid && err == nil {
			t.Errorf("assertValidMatrix() did not return an error for an invalid matrix")
			return
		}
	}
}

func TestAdjustContrast(t *testing.T) {

	width, height := rand.Intn(3841), rand.Intn(2161) // Random dimensions up to 4k
	randomImage := generateRandomImage(width, height)

	// Example values for m (contrast factor) and b (brightness offset)
	m := 1.2 // Contrast factor
	b := 0.0 // Brightness offset

	// Apply the contrast adjustment
	adjustedImage, err := manipulations.AdjustContrast(randomImage, m, b)

	// Check if any error is returned
	if err != nil {
		t.Errorf("AdjustContrast() returned an error: %v", err)
	}

	// Assert that the adjusted image is valid
	if err := assertValidMatrix(adjustedImage); err != nil {
		t.Errorf("AdjustContrast() resulted in an invalid matrix: %v", err)
	}
}

func TestAdjustLuminosity(t *testing.T) {

	width, height := rand.Intn(3841), rand.Intn(2161) // Random dimensions up to 4k
	randomImage := generateRandomImage(width, height)

	b := 50.0 // Brightness offset

	// Apply the contrast adjustment
	adjustedImage, err := manipulations.AdjustLuminosity(randomImage, b)

	// Check if any error is returned
	if err != nil {
		t.Errorf("AdjustLuminosity() returned an error: %v", err)
	}

	// Assert that the adjusted image is valid
	if err := assertValidMatrix(adjustedImage); err != nil {
		t.Errorf("AdjustLuminosity() resulted in an invalid matrix: %v", err)
	}
}
