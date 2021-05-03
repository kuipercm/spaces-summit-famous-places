package web

import (
	"embed"
	"io/fs"
	"net/http"
)

// https://golang.org/pkg/embed/
var (
	//go:embed static
	statics embed.FS
	//go:embed static/index.html
	indexHtml []byte
)

func NewHandler() (http.Handler, error) {
	// Make sure that this filesystem starts at /
	fsys, err := fs.Sub(statics, "static")
	if err != nil {
		return nil, err
	}

	return http.FileServer(http.FS(fsys)), nil
}
