package wall

import (
	"log"
	"sync"
	"time"
)

// Processor is used to transform photos in a pipeline like resizing and checksumming
type Processor interface {
	Process(p Photo) (Photo, error)
}

// ProcessorFunc can be used instead of the Processor type
type ProcessorFunc func(p Photo) (Photo, error)

// Process transforms ProcessorFunc to a Processor
func (f ProcessorFunc) Process(p Photo) (Photo, error) {
	return f(p)
}

// Observer gets notified if something changes on the wall
type Observer func(p Photo)

// Photowall represents a wall of photos
type Photowall interface {
	AddPhotoFromFile(name string, createdAt time.Time) error
	AddPhoto(p Photo) error
	RemovePhoto(photo Photo)
	OnAdd(o Observer)
	OnRemove(o Observer)
	Photos() Photos
}

// Wall represents a collection of photos, create with Create
type Wall struct {
	processors      []Processor
	photos          Photos
	mutexPhotos     *sync.RWMutex
	listenersAdd    []Observer
	listenersRemove []Observer
}

// Create a new photowall
//
// Processors should be set after creation
// 	wall.SetProcessors([]Processor{
//			NewResizer(1920, 1080),
//			NewStore("./storage"),
// 	})
func Create() *Wall {
	return &Wall{
		mutexPhotos: &sync.RWMutex{},
		processors: []Processor{
			NewResizer(1920, 1080),
			NewStore("./storage"),
		},
	}
}

// SetProcessors sets the list of registered processors
func (w *Wall) SetProcessors(ps []Processor) {
	w.processors = ps
}

// Processors returns the list of registered processors
func (w *Wall) Processors() []Processor {
	return w.processors
}

// AddPhotoFromFile adds a new photo to the wall
func (w *Wall) AddPhotoFromFile(name string, createdAt time.Time) error {
	p := NewPhoto(name, 0, 0, "", createdAt)
	return w.process(p)
}

// AddPhoto adds a new photo to the wall
func (w *Wall) AddPhoto(p Photo) error {
	return w.process(p)
}

func (w *Wall) storePhoto(p Photo) {
	log.Printf("Store photo: %s", p.Name())
	w.mutexPhotos.Lock()
	w.photos = append(w.photos, p)
	w.mutexPhotos.Unlock()
	w.notifyAdd(p)
}

func (w *Wall) process(photo Photo) error {
	var err error
	for _, p := range w.processors {
		photo, err = p.Process(photo)
		if err != nil {
			// TODO: handle error
			return err
		}
	}

	w.storePhoto(photo)
	return nil
}

// RemovePhoto removes a photo from the wall
func (w *Wall) RemovePhoto(photo Photo) {
	w.mutexPhotos.Lock()
	var index = -1
	for i, p := range w.photos {
		if p == photo {
			index = i
			break
		}
	}
	if index != -1 {
		w.photos = append(w.photos[:index], w.photos[index+1:]...)
	}
	w.mutexPhotos.Unlock()
	w.notifyRemove(photo)
}

func (w *Wall) notifyAdd(p Photo) {
	for _, o := range w.listenersAdd {
		go o(p)
	}
}

func (w *Wall) notifyRemove(p Photo) {
	for _, o := range w.listenersRemove {
		go o(p)
	}
}

// OnAdd registers an Observer which is called when a photo was added to the wall
func (w *Wall) OnAdd(o Observer) {
	w.listenersAdd = append(w.listenersAdd, o)
}

// OnRemove registers an Observer which is called when a photo was removed to the wall
func (w *Wall) OnRemove(o Observer) {
	w.listenersRemove = append(w.listenersRemove, o)
}

// Photos returns all photos on the wall
func (w Wall) Photos() Photos {
	w.mutexPhotos.RLock()
	b := make([]Photo, len(w.photos))
	copy(b, w.photos)
	w.mutexPhotos.RUnlock()
	return b
}
