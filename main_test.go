package matrix_image_manipulation

import (
	"image"
	"os"
	"testing"
)

func loadImage(filePath string) (image.Image, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	return img, err
}

// TestMatrixContinuity tests that if an image is converted to a matrix, and then back, it remains the same
func TestMatrixContinuity(t *testing.T) {
	// Load test images
	image, err := loadImage("test_images/image1.png")
	if err != nil {
		t.Fatalf("Failed to load image1: %s", err)
	}

}
