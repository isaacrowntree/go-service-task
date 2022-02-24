package main

import (
	"log"
	"net/http"

	"github.com/isaacrowntree/go-service-task/pkg/handler"
)

func main() {
	http.Handle("/", handler.Handler{H: handler.GetResults})
	log.Println("Listening for requests at http://localhost:8000/")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
