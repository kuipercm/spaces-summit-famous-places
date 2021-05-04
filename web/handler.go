package web

import (
	"embed"
	"io/fs"
	"net/http"
)

// https://golang.org/pkg/embed/
//go:embed static
var statics embed.FS

func NewHandler() (http.Handler, error) {
	// Make sure that this filesystem starts at /
	fsys, err := fs.Sub(statics, "static")
	if err != nil {
		return nil, err
	}

	return http.FileServer(http.FS(fsys)), nil
}
