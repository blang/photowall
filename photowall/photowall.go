package photowall

import (
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

// Photowall represents a collection of photos, create with Create
type Photowall struct {
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
func Create() *Photowall {
	return &Photowall{
		mutexPhotos: &sync.RWMutex{},
		processors: []Processor{
			NewResizer(1920, 1080),
			NewStore("./storage"),
		},
	}
}

// SetProcessors sets the list of registered processors
func (w *Photowall) SetProcessors(ps []Processor) {
	w.processors = ps
}

// Processors returns the list of registered processors
func (w *Photowall) Processors() []Processor {
	return w.processors
}

// AddPhotoFromFile adds a new photo to the wall
func (w *Photowall) AddPhotoFromFile(name string, createdAt time.Time) {
	p := NewPhoto(name, 0, 0, "", createdAt)
	go w.process(p)
}

// AddPhoto adds a new photo to the wall
func (w *Photowall) AddPhoto(p Photo) {
	go w.process(p)
}

func (w *Photowall) storePhoto(p Photo) {
	w.mutexPhotos.Lock()
	w.photos = append(w.photos, p)
	w.mutexPhotos.Unlock()
	w.notifyAdd(p)
}

func (w *Photowall) process(photo Photo) {
	var err error
	for _, p := range w.processors {
		photo, err = p.Process(photo)
		if err != nil {
			// handle error
			break
		}
	}
	if err == nil {
		w.storePhoto(photo)
	}
}

// RemovePhoto removes a photo from the wall
func (w *Photowall) RemovePhoto(photo Photo) {
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

func (w *Photowall) notifyAdd(p Photo) {
	for _, o := range w.listenersAdd {
		go o(p)
	}
}

func (w *Photowall) notifyRemove(p Photo) {
	for _, o := range w.listenersRemove {
		go o(p)
	}
}

// OnAdd registers an Observer which is called when a photo was added to the wall
func (w *Photowall) OnAdd(o Observer) {
	w.listenersAdd = append(w.listenersAdd, o)
}

// OnRemove registers an Observer which is called when a photo was removed to the wall
func (w *Photowall) OnRemove(o Observer) {
	w.listenersRemove = append(w.listenersRemove, o)
}

// Photos returns all photos on the wall
func (w Photowall) Photos() Photos {
	w.mutexPhotos.RLock()
	b := make([]Photo, len(w.photos))
	copy(b, w.photos)
	w.mutexPhotos.RUnlock()
	return b
}
