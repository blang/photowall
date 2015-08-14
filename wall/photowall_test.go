package wall

import (
	"image"
	"image/jpeg"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func createTestImg() (string, error) {
	m := image.NewRGBA(image.Rect(0, 0, 1000, 2000))
	fout, err := ioutil.TempFile("", "imagetest")
	if err != nil {
		return "", err
	}
	defer fout.Close()
	err = jpeg.Encode(fout, m, nil)
	if err != nil {
		return "", err
	}
	return fout.Name(), nil
}

func TestPhotowall(t *testing.T) {
	w := Create()
	var (
		proc1        Photo
		proc2        Photo
		addCalled    Photo
		removeCalled Photo
	)
	w.SetProcessors([]Processor{
		ProcessorFunc(func(p Photo) (Photo, error) {
			proc1 = p
			return p, nil
		}),
		ProcessorFunc(func(p Photo) (Photo, error) {
			proc2 = p
			return p, nil
		}),
	})
	w.OnAdd(Observer(func(p Photo) {
		addCalled = p
	}))

	w.OnRemove(Observer(func(p Photo) {
		removeCalled = p
	}))

	imgName, err := createTestImg()
	if err != nil {
		t.Fatalf("Could not create test image: %s", err)
	}
	defer os.Remove(imgName)
	err = w.AddPhotoFromFile(imgName, time.Now())
	if err != nil {
		t.Errorf("Error adding photo: %s", err)
	}

	// Check processors
	if p := proc1; p == nil || p.Name() != imgName {
		t.Errorf("Processor1 received wrong image: %s", p)
	}

	if p := proc2; p == nil || p.Name() != imgName {
		t.Errorf("Processor2 received wrong image: %s", p)
	}

	// Check add observer
	if p := addCalled; p == nil || p.Name() != imgName {
		t.Errorf("Add received wrong image: %s", p)
	}

	// Check photos
	photos := w.Photos()
	if len(photos) != 1 {
		t.Errorf("Invalid amount of photos: %s", photos)
	}

	w.RemovePhoto(photos[0])

	if p := removeCalled; p == nil || p.Name() != imgName {
		t.Errorf("Remove received wrong image: %s", p)
	}

	photos = w.Photos()
	if len(photos) != 0 {
		t.Errorf("Invalid amount of photos: %s", photos)
	}
}
