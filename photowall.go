package main

import (
	"flag"
	"github.com/blang/photowall/wall"
	"github.com/blang/photowall/web"
	"io/ioutil"
	"log"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

var listen = flag.String("listen", ":8000", "Listen addr")
var storeDir = flag.String("storedir", "./imgs", "Storage directory")
var argAllowedExts = flag.String("allow", "png,jpg", "Allowed file extensions")
var argImgWidth = flag.Uint("img_width", 1920, "Resize bigger images to this width")
var argImgHeight = flag.Uint("img_height", 1080, "Resize bigger images to this height")
var argMaxFileSize = flag.Int("filesize_max", 10, "Maximum upload filesize in MB")

func baseDir() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Dir(file)
}

func main() {
	flag.Parse()
	pwall := wall.Create()
	pwall.SetProcessors([]wall.Processor{
		wall.Importer(),
	})
	// Restore existing images using Importer Processor
	restoreFromDirectory(pwall, filepath.Join(baseDir(), *storeDir))

	// Set Production processors
	pwall.SetProcessors([]wall.Processor{
		wall.NewResizer(*argImgWidth, *argImgHeight),
		wall.NewStore(filepath.Join(baseDir(), *storeDir)),
	})

	server := web.NewServer(pwall, filepath.Join(baseDir(), "/static"), filepath.Join(baseDir(), *storeDir), int64(*argMaxFileSize)*1024*1025, *argAllowedExts)
	log.Fatal(server.Run(*listen))
}

func restoreFromDirectory(wall wall.Photowall, path string) {
	log.Printf("Restore store from directory: %s\n", path)
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Printf("Error reading directory: %s", err)
		return
	}
	wg := &sync.WaitGroup{}

	for _, f := range files {
		if strings.ToLower(filepath.Ext(f.Name())) == ".jpg" {
			fullpath := filepath.Join(path, f.Name())
			wg.Add(1)
			go func(path string) {
				err := wall.AddPhotoFromFile(path, time.Now())
				if err == nil {
					log.Printf("Added file: %s", fullpath)
				} else {
					log.Printf("Error adding file %s: %s", fullpath, err)
				}
				wg.Done()
			}(fullpath)
		}
	}
	wg.Wait()
}
