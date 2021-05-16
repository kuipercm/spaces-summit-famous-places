package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/kuipercm/spaces-summit-famous-places/firestore"
)

type fileHandler struct {
	fireStore firestore.Store
}

func newFileHandler(f firestore.Store) fileHandler {
	return fileHandler{
		fireStore: f,
	}
}

func (m fileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	qLimit := r.URL.Query().Get("limit")
	qOffset := r.URL.Query().Get("offset")
	if qLimit == "" || qOffset == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("You should specify both limit and offset. offset = " + qOffset + ", qLimit = " + qLimit))
		return
	}

	limit, err := strconv.Atoi(qLimit)
	if err != nil {
		fmt.Printf("parse qLimit string to int %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	offset, err := strconv.Atoi(qOffset)
	if err != nil {
		fmt.Printf("parse qOffset string to int %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := m.fireStore.List(r.Context(), limit, offset)
	if err != nil {
		fmt.Printf("firestore::list %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(res)
}
