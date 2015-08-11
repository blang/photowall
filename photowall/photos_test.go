package photowall

import (
	"testing"
	"time"
)

func TestPhoto(t *testing.T) {
	createdAt := time.Now()
	p := NewPhoto("test", 100, 200, "png", createdAt)
	if p.Name() != "test" {
		t.Errorf("Wrong name: %s", p.Name())
	}
	size := p.Bounds().Size()
	if size.X != 100 {
		t.Errorf("Wrong size X: %d", size.X)
	}
	if size.Y != 200 {
		t.Errorf("Wrong size Y: %d", size.Y)
	}
	if p.Format() != "png" {
		t.Errorf("Wrong format: %s", p.Format())
	}
	if p.CreatedAt() != createdAt {
		t.Errorf("Wrong createdAt: %s", p.CreatedAt())
	}
}

func TestSort(t *testing.T) {
	createdAt := time.Now()
	ps := []Photo{
		NewPhoto("b", 100, 200, "png", createdAt.Add(5*time.Second)),
		NewPhoto("a", 100, 200, "png", createdAt),
	}

	SortPhotoSlice(ps)
	if first := ps[0].Name(); first != "a" {
		t.Fatalf("Wrong sorting, first item is: %s", first)
	}

	photos := Photos(ps)
	photos.Swap(0, 1)
	SortPhotos(photos)
	if first := photos[0].Name(); first != "a" {
		t.Errorf("Wrong sorting, first item is: %s", first)
	}

}
