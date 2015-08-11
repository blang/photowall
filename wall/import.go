package wall

import (
	"bufio"
	"image"
	_ "image/jpeg"
	"os"
)

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
	img, _, err := image.Decode(imgReader)
	if err != nil {
		return nil, err
	}

	dims := img.Bounds().Size()

	return NewPhoto(p.Name(), dims.X, dims.Y, "jpg", p.CreatedAt()), nil
}
