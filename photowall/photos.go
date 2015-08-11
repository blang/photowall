package photowall

import (
	"image"
	"sort"
	"time"
)

type wallPhoto struct {
	name      string
	bounds    image.Rectangle
	format    string
	createdAt time.Time
}

// NewPhoto creates a new photo
func NewPhoto(name string, width, height int, format string, createdAt time.Time) Photo {
	return wallPhoto{
		name:      name,
		bounds:    image.Rect(0, 0, width, height),
		format:    format,
		createdAt: createdAt,
	}
}

func (p wallPhoto) Name() string {
	return p.name
}

func (p wallPhoto) Bounds() image.Rectangle {
	return p.bounds
}

func (p wallPhoto) Format() string {
	return p.format
}

func (p wallPhoto) CreatedAt() time.Time {
	return p.createdAt
}

// Photo represents a photo on the photowall
type Photo interface {
	Name() string
	Bounds() image.Rectangle // width, height
	Format() string          // png, jpeg
	CreatedAt() time.Time
}

// Photos is a collection of photos
type Photos []Photo

// Len returns length of photo collection
func (s Photos) Len() int {
	return len(s)
}

// Swap swaps two photos inside the collection by its indices
func (s Photos) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Less checks if photo at index i is created before photo at index j
func (s Photos) Less(i, j int) bool {
	return s[i].CreatedAt().Before(s[j].CreatedAt())
}

// SortPhotos sorts photos
func SortPhotos(photos Photos) {
	sort.Sort(photos)
}

// SortPhotoSlice sorts a slice of photos
func SortPhotoSlice(photos []Photo) {
	sort.Sort(Photos(photos))
}
