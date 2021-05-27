package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/kuipercm/spaces-summit-famous-places/firestore"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type fileHandler struct {
	fireStore firestore.Store
	tracer    trace.Tracer
}

func newFileHandler(f firestore.Store) fileHandler {
	t := otel.Tracer("ssfp/uploads")

	return fileHandler{
		fireStore: f,
		tracer:    t,
	}
}

func (m fileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, span := m.tracer.Start(r.Context(), "http/list")
	defer span.End()

	qLastCreationDate := r.URL.Query().Get("creationDate")
	if qLastCreationDate != "" {
		m.byLastCreationDate(ctx, w, r, qLastCreationDate)
		return
	}

	qLimit := r.URL.Query().Get("limit")
	qOffset := r.URL.Query().Get("offset")
	if qLimit == "" || qOffset == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("If not specifing by creationDate, you should specify both limit and offset. offset = " + qOffset + ", qLimit = " + qLimit))
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

	fsCtx, fsSpan := m.tracer.Start(ctx, "firestore/list")
	res, err := m.fireStore.List(fsCtx, limit, offset)
	fsSpan.End()

	if err != nil {
		fmt.Printf("firestore::list %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(res)
}

func (m fileHandler) byLastCreationDate(ctx context.Context, w http.ResponseWriter, r *http.Request, qLastCreationDate string) {
	lastCreationDateMillis, err := strconv.ParseInt(qLastCreationDate, 10, 64)
	if err != nil {
		fmt.Printf("parse qLastCreationDate string to int %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	lastCreationDate := time.Unix(0, lastCreationDateMillis*int64(time.Millisecond))

	ctx, span := m.tracer.Start(ctx, "firestore/list/ByCreationDate")
	res, err := m.fireStore.ListByCreationDate(ctx, lastCreationDate)
	defer span.End()

	if err != nil {
		fmt.Printf("firestore::ListByCreationDate %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(res)
}
