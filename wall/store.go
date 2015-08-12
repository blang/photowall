package wall

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"io"
	"os"
	"path/filepath"
)

// Store processes photos, stores them inside a given directory and checks for duplicates
type Store struct {
	dir    string
	chsums map[string]struct{}
	namer  Namer
}

// NewStore creates a new store processor
func NewStore(directory string) *Store {
	return &Store{
		dir:    directory,
		chsums: make(map[string]struct{}),
		namer:  NewDateNamer("2006-01-02_150405"),
	}
}

// SetNamer sets the Namer for filenames
func (s *Store) SetNamer(namer Namer) {
	s.namer = namer
}

// Process copy the photo to the store directory and discard it if it's a dup
func (s *Store) Process(p Photo) (Photo, error) {
	newBaseName := s.namer.Name(p) + "." + p.Format()
	fin, err := os.Open(p.Name())
	if err != nil {
		return nil, err
	}
	defer func() {
		fin.Close()
		os.Remove(fin.Name())
	}()

	newName := filepath.Join(s.dir, newBaseName)
	fout, err := os.Create(newName)
	if err != nil {
		return nil, err
	}
	defer fout.Close()
	hash := sha1.New()
	imgReader := io.TeeReader(fin, hash)
	_, err = io.Copy(fout, imgReader)
	if err != nil {
		return nil, err
	}

	chsum := hex.EncodeToString(hash.Sum(nil))
	if _, ok := s.chsums[chsum]; ok {
		//TODO: Remove
		return nil, errors.New("File already exists")
	}
	return NewPhoto(newName, p.Bounds().Size().X, p.Bounds().Size().Y, p.Format(), p.CreatedAt()), nil
}
