package web

import (
	"github.com/blang/photowall/wall"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

// Server represents a http server serving the photowall and upload functionality
type Server struct {
	*gin.Engine
	wall            wall.Photowall
	maxSize         int64
	validExtensions map[string]struct{}
	storageDir      string
}

func buildValidExtensions(extensions string) map[string]struct{} {
	extMap := make(map[string]struct{})
	exts := strings.Split(extensions, ",")
	for _, extstr := range exts {
		extstr = strings.ToLower(strings.TrimSpace(extstr))
		extstr = strings.Replace(extstr, ".", "", -1)
		if extstr != "" {
			extMap[extstr] = struct{}{}
		}
	}
	return extMap
}

func (s Server) validExtension(name string) (string, bool) {
	ext := filepath.Ext(name)
	if ext == "" {
		return "", false
	}
	ext = strings.ToLower(strings.TrimSpace(ext))
	ext = strings.Replace(ext, ".", "", -1)
	if _, ok := s.validExtensions[ext]; ok {
		return ext, true
	}
	return "", false

}

// NewServer creates a new Server instance
func NewServer(wall wall.Photowall, staticDir string, storageDir string, maxSize int64, validExtensions string) *Server {
	s := &Server{}
	s.wall = wall
	s.maxSize = maxSize
	s.storageDir = storageDir
	s.validExtensions = buildValidExtensions(validExtensions)

	router := gin.Default()
	router.LoadHTMLGlob(filepath.Join(staticDir, "/templates/*"))

	router.Static("/imgs", storageDir)
	router.Static("/assets", filepath.Join(staticDir, "/assets"))
	router.StaticFile("/wall", filepath.Join(staticDir, "/wall.html"))
	router.StaticFile("/admin", filepath.Join(staticDir, "/admin.html"))
	router.StaticFile("/success", filepath.Join(staticDir, "/success.html"))
	router.StaticFile("/", filepath.Join(staticDir, "/upload.html"))
	router.POST("/api/upload", s.handleUpload)
	router.GET("/api/wall", s.handleAPIWall)
	s.Engine = router
	return s
}

type exportPhoto struct {
	Name      string `json:"name"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	CreatedAt string `json:"created_at"`
}

func exportPhotos(ps wall.Photos) []exportPhoto {
	var export []exportPhoto
	wall.SortPhotos(ps)
	for _, p := range ps {
		export = append(export, exportPhoto{
			Name:      filepath.Base(p.Name()),
			Width:     p.Bounds().Size().X,
			Height:    p.Bounds().Size().Y,
			CreatedAt: p.CreatedAt().String(),
		})
	}
	return export
}

func (s Server) handleAPIWall(c *gin.Context) {
	c.JSON(http.StatusOK, exportPhotos(s.wall.Photos()))
}

func (s Server) handleUpload(c *gin.Context) {
	if c.Request.ContentLength > s.maxSize {
		http.Error(c.Writer, "request too large", http.StatusExpectationFailed)
		return
	}
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, s.maxSize)
	err := c.Request.ParseMultipartForm(1024)
	if err != nil {
		log.Printf("Could not get file from form: %s\n", err)
		c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{})
		return
	}
	file, handler, err := c.Request.FormFile("pic")
	if err != nil {
		log.Printf("Could not get file from form: %s\n", err)
		c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{})
		return
	}
	defer file.Close()
	ext, ok := s.validExtension(handler.Filename)
	if !ok {
		log.Printf("Invalid file extension: %s", handler.Filename)
		c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{})
		return
	}
	f, err := ioutil.TempFile("", ext)
	if err != nil {
		log.Printf("Could not create file: %s\n", err)
		c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{})
		return
	}
	defer f.Close()
	_, err = io.Copy(f, file)
	if err != nil {
		log.Printf("File error: %s\n", err)
		c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{})
		return
	}

	s.wall.AddPhotoFromFile(f.Name(), time.Now())

	http.Redirect(c.Writer, c.Request, "/success", http.StatusFound)
}
