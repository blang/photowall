package wall

import (
	"bufio"
	"errors"
	"image"
	_ "image/jpeg" // Support jpeg image format
	"os"
)

// Importer creates a processor, checking the file for valid image
func Importer() Processor {
	return ProcessorFunc(importProcess)
}

func importProcess(p Photo) (Photo, error) {
	file, err := os.Open(p.Name())
	if err != nil {
		return nil, err
	}
	defer file.Close()
	imgReader := bufio.NewReader(file)

	// decode jpeg into image.Image
	img, format, err := image.Decode(imgReader)
	if err != nil {
		return nil, err
	}
	if format != "jpeg" {
		return nil, errors.New("Not a valid jpeg")
	}

	dims := img.Bounds().Size()

	return NewPhoto(p.Name(), dims.X, dims.Y, "jpg", p.CreatedAt()), nil
}
