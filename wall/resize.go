package wall

import (
	"bufio"
	"github.com/nfnt/resize"
	"image"
	_ "image/gif" // Support gif format
	"image/jpeg"
	_ "image/png" // Support png format
	"io/ioutil"
	"os"
)

// Resizer resizes photos
type Resizer struct {
	MaxWidth  uint
	MaxHeight uint
}

// NewResizer creates a new instance of Resizer
func NewResizer(maxWidth, maxHeight uint) Resizer {
	return Resizer{
		MaxWidth:  maxWidth,
		MaxHeight: maxHeight,
	}
}

// Process starts the resizing of the photo
func (r Resizer) Process(p Photo) (Photo, error) {
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

	// resize to width 1000 using Lanczos resampling
	// and preserve aspect ratio
	dims := img.Bounds().Size()
	var newDimX uint
	var newDimY uint
	if dims.X > dims.Y {
		newDimX = r.MaxWidth
	} else {
		newDimY = r.MaxHeight
	}

	var m image.Image
	if uint(dims.X) > r.MaxWidth || uint(dims.Y) > r.MaxHeight {
		m = resize.Resize(newDimX, newDimY, img, resize.Lanczos3)
	} else {
		// Don't resize small images
		m = img
	}
	out, err := ioutil.TempFile("", ".jpg")
	if err != nil {
		return nil, err
	}
	defer out.Close()
	// write new image to file
	jpeg.Encode(out, m, nil)
	newdims := m.Bounds().Size()
	os.Remove(p.Name())
	return NewPhoto(out.Name(), newdims.X, newdims.Y, "jpg", p.CreatedAt()), nil
}
