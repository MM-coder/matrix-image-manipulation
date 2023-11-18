package matrix_image_manipulation

import (
	"fmt"
	"image"
	"image/color"
	"math/rand"
	"os"
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

// TestMatrixContinuity tests that if an image is converted to a matrix, and then back, it remains the same
func TestMatrixContinuity(t *testing.T) {

	var testImagePath string = ".github/test_images/gnome.png"

	matrix, err := readImageToMatrix(testImagePath)
	if err != nil {
		t.Fatalf("Failed to load the test image: %s", err)
		return
	}

	temporaryPath := t.TempDir() + "gnome.png"

	_ = writeImageFromMatrix(matrix, temporaryPath)

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
	width, height := rand.Intn(100)+1, rand.Intn(100)+1 // As  rand.Intn returns from [0, n] we must add one to assert n is never 0
	matrix := Make2D[[4]int32](width, height)

	// Check if the number of rows is n
	if len(matrix) != width {
		t.Errorf("Expected %d rows, got %d", width, len(matrix))
	}

	// Check if each row has m elements
	for i, row := range matrix {
		if len(row) != height {
			t.Errorf("Expected %d elements in row %d, got %d", height, i, len(row))
		}
	}

	// Check if each element's capacity is 4, which is as requirement
	for _, row := range matrix {
		for _, elem := range row {
			if !(cap(elem) == 4) {
				t.Errorf("Expected element to be initialized to an empty slice, got %v", elem)
			}
		}
	}
}
