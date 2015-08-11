package wall

import (
	"image"
	_ "image/jpeg"
	"image/png"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func createResizeTestImg(w, h int) (string, error) {
	m := image.NewRGBA(image.Rect(0, 0, w, h))
	fout, err := ioutil.TempFile("", "imagetest")
	if err != nil {
		return "", err
	}
	defer fout.Close()
	err = png.Encode(fout, m)
	if err != nil {
		return "", err
	}
	return fout.Name(), nil
}

func TestResize(t *testing.T) {
	fo, err := createResizeTestImg(1000, 2000)
	if err != nil {
		t.Fatal("Could not create tmp file")
	}

	defer os.Remove(fo)

	var resizer Processor
	resizer = NewResizer(100, 200)

	p := NewPhoto(fo, 0, 0, "", time.Now())

	newP, err := resizer.Process(p)
	if err != nil {
		t.Fatalf("Error resizing img: %s", err)
	}
	size := newP.Bounds().Size()
	if size.X != 100 {
		t.Errorf("Wrong width: %d", size.X)
	}
	if size.Y != 200 {
		t.Errorf("Wrong height: %d", size.Y)
	}

	fi, err := os.Open(newP.Name())
	if err != nil {
		t.Errorf("Error opening new file: %s", err)
	}
	defer fi.Close()
	newimg, format, err := image.Decode(fi)
	if err != nil {
		t.Errorf("Could not decode image: %s", err)
	}
	if format != "jpeg" {
		t.Errorf("Wrong format: %s", format)
	}
	imgsize := newimg.Bounds().Size()
	if size.X != imgsize.X {
		t.Errorf("Image has wrong width: %d", imgsize.X)
	}

	if size.Y != imgsize.Y {
		t.Errorf("Image has wrong height: %d", imgsize.Y)
	}

	// Make sure tmp file is removed
	_, err = os.Stat(fo)
	if err == nil {
		t.Errorf("Input photo file was not removed")
	}
}

func TestResizeTooSmall(t *testing.T) {
	m := image.NewRGBA(image.Rect(0, 0, 50, 50))
	fo, err := ioutil.TempFile("", "imagetest")
	if err != nil {
		t.Fatal("Could not create tmp file")
	}

	defer fo.Close()
	png.Encode(fo, m)

	var resizer Processor
	resizer = NewResizer(100, 200)

	p := NewPhoto(fo.Name(), 0, 0, "", time.Now())

	newP, err := resizer.Process(p)
	if err != nil {
		t.Fatalf("Error resizing img: %s", err)
	}
	size := newP.Bounds().Size()
	if size.X != 50 {
		t.Errorf("Wrong width: %d", size.X)
	}
	if size.Y != 50 {
		t.Errorf("Wrong height: %d", size.Y)
	}

	fi, err := os.Open(newP.Name())
	if err != nil {
		t.Errorf("Error opening new file: %s", err)
	}
	defer fi.Close()
	newimg, format, err := image.Decode(fi)
	if err != nil {
		t.Errorf("Could not decode image: %s", err)
	}
	if format != "jpeg" {
		t.Errorf("Wrong format: %s", format)
	}
	imgsize := newimg.Bounds().Size()
	if size.X != imgsize.X {
		t.Errorf("Image has wrong width: %d", imgsize.X)
	}

	if size.Y != imgsize.Y {
		t.Errorf("Image has wrong height: %d", imgsize.Y)
	}

}
