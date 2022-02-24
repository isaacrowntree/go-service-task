package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/isaacrowntree/go-service-task/pkg/parser"
	"github.com/isaacrowntree/go-service-task/pkg/reader"
	"github.com/isaacrowntree/go-service-task/pkg/slicer"
	"github.com/isaacrowntree/go-service-task/pkg/structs"
)

type Error interface {
	error
	Status() int
}

type StatusError struct {
	Code int
	Err  error
}

func (se StatusError) Error() string {
	return se.Err.Error()
}

func (se StatusError) Status() int {
	return se.Code
}

type Handler struct {
	H func(w http.ResponseWriter, r *http.Request) error
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h.H(w, r)
	if err != nil {
		switch e := err.(type) {
		case Error:
			log.Printf("HTTP %d - %s", e.Status(), e)
			http.Error(w, e.Error(), e.Status())
		default:
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}
	}
}

func GetResults(w http.ResponseWriter, r *http.Request) error {
	decoder := json.NewDecoder(r.Body)
	var p structs.Parameters
	err := decoder.Decode(&p)
	if err != nil {
		return StatusError{404, err}
	}

	log.Print("Beginning search...")
	start := time.Now()

	files := reader.GetFileList(p.Filename)
	chunks := slicer.ChunkSlice(files, runtime.GOMAXPROCS(0))
	values := parser.ParseFile(chunks, p)

	elapsed := time.Since(start)
	log.Printf("Finished! Log parser took %s for %s", elapsed, p.Filename)

	j, err := json.Marshal(values)
	if err != nil {
		return StatusError{500, err}
	}

	w.Write(j)
	return nil
}
