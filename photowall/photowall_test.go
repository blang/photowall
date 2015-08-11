package photowall

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
	processor1Called := make(chan Photo, 1)
	processor2Called := make(chan Photo, 1)
	addCalled := make(chan Photo, 1)
	removeCalled := make(chan Photo, 1)
	w.SetProcessors([]Processor{
		ProcessorFunc(func(p Photo) (Photo, error) {
			processor1Called <- p
			return p, nil
		}),
		ProcessorFunc(func(p Photo) (Photo, error) {
			processor2Called <- p
			return p, nil
		}),
	})
	w.OnAdd(Observer(func(p Photo) {
		addCalled <- p
	}))

	w.OnRemove(Observer(func(p Photo) {
		removeCalled <- p
	}))

	imgName, err := createTestImg()
	if err != nil {
		t.Fatalf("Could not create test image: %s", err)
	}
	defer os.Remove(imgName)
	w.AddPhotoFromFile(imgName, time.Now())

	// Check processor
	select {
	case <-time.After(3 * time.Second):
		t.Errorf("Processor1 was not called in time")
	case p := <-processor1Called:
		if p.Name() != imgName {
			t.Errorf("Processor1 received wrong image: %s", p)
		}
	}

	select {
	case <-time.After(3 * time.Second):
		t.Errorf("Processor2 was not called in time")
	case p := <-processor2Called:
		if p.Name() != imgName {
			t.Errorf("Processor2 received wrong image: %s", p)
		}
	}

	// Check add observer
	select {
	case <-time.After(3 * time.Second):
		t.Errorf("Add was not called in time")
	case p := <-addCalled:
		if p.Name() != imgName {
			t.Errorf("Add received wrong image: %s", p)
		}
	}

	// Check photos
	photos := w.Photos()
	if len(photos) != 1 {
		t.Errorf("Invalid amount of photos: %s", photos)
	}

	w.RemovePhoto(photos[0])

	select {
	case <-time.After(3 * time.Second):
		t.Errorf("Remove was not called in time")
	case p := <-removeCalled:
		if p.Name() != imgName {
			t.Errorf("Remove received wrong image: %s", p)
		}
	}

	photos = w.Photos()
	if len(photos) != 0 {
		t.Errorf("Invalid amount of photos: %s", photos)
	}
}
