package handlers

import (
	"github.com/google/uuid"
	"github.com/kuipercm/spaces-summit-famous-places/storage"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

type SpaHandler struct {
	StaticPath string
	IndexPath  string
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
	if os.IsNotExist(err) {
		http.ServeFile(w, r, filepath.Join(h.StaticPath, h.IndexPath))
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// otherwise, use http.FileServer to serve the static dir
	http.FileServer(http.Dir(h.StaticPath)).ServeHTTP(w, r)
}

type MultipartUploadHandler struct {
	Storage storage.GcpStorage
	MaxSize int64
}

func (m MultipartUploadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(m.MaxSize) // maxMemory 32MB = 32 << 20
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for _, h := range r.MultipartForm.File["photo"] {
		fileId, err := uuid.NewRandom()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		//err = m.saveToDisk(h, fileId.String() + filepath.Ext(h.Filename))
		fileName := fileId.String() + filepath.Ext(h.Filename)
		err = m.Storage.StoreFile(fileName, h)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	w.WriteHeader(200)
	return
}

func (m MultipartUploadHandler) saveToDisk(h *multipart.FileHeader, fileName string) error {
	file, err := h.Open()
	if err != nil {
		return err
	}

	outputFile, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	readBytes := make([]byte, 64)
	_, err = file.Read(readBytes)
	if err != nil {
		return err
	}

	for len(readBytes) > 0 {
		_, err := outputFile.Write(readBytes)
		if err != nil {
			return err
		}

		_, err = file.Read(readBytes)
		if err != nil {
			return err
		}
	}

	return nil
}
