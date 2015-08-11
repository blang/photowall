package photowall

import (
	"image"
	"image/jpeg"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func createStoreTestImg() (string, error) {
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

func TestStore(t *testing.T) {
	dirName, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatalf("Could not create tmp dir: %s", err)
	}
	s := NewStore(dirName)
	pName, err := createStoreTestImg()
	if err != nil {
		t.Fatalf("Could not test image: %s", err)
	}
	defer os.Remove(pName)

	// Read image back
	statIn, err := os.Stat(pName)
	if err != nil {
		t.Fatalf("Could not stat tmp file: %s", err)
	}

	const extension = "jpg"
	inPhoto := NewPhoto(pName, 0, 0, extension, time.Now())
	outPhoto, err := s.Process(inPhoto)
	if err != nil {
		t.Errorf("Error while processing: %s", err)
	}
	if filepath.Dir(outPhoto.Name()) != dirName {
		t.Errorf("Image stored in wrong directory: %s", outPhoto.Name())
	}
	if ext := filepath.Ext(outPhoto.Name()); ext != "."+extension {
		t.Errorf("Image stored has wrong extension: %s", ext)
	}

	statOut, err := os.Stat(outPhoto.Name())
	if err != nil {
		t.Fatalf("Could not stat stored file: %s", err)
	}

	if statIn.Size() != statOut.Size() {
		t.Errorf("Stored file has wrong size: %d", statOut.Size())
	}

	// Make sure tmp file is removed
	_, err = os.Stat(pName)
	if err == nil {
		t.Errorf("Input photo file was not removed")
	}
}
