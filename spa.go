package main

import (
	"net/http"
	"os"
	"path/filepath"
)

type SpaHandler struct {
	StaticPath string
	IndexPath  string
}

// We can use go:embed here since 1.16
// https://golang.org/pkg/embed/
// This makes it one deployable binary instead of having to ship the web files as well.
func newSpaHandler(staticPath, indexPath string) SpaHandler {
	return SpaHandler{
		StaticPath: staticPath,
		IndexPath:  filepath.Join(staticPath, indexPath),
	}
}

func (h SpaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// get the absolute path to prevent directory traversal
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// prepend the path with the path to the static directory
	path = filepath.Join(h.StaticPath, path)

	// check whether a file exists at the given path
	_, err = os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			http.ServeFile(w, r, h.IndexPath)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// otherwise, use http.FileServer to serve the static dir
	http.FileServer(http.Dir(h.StaticPath)).ServeHTTP(w, r)
}
