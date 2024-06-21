package main

import (
	"log"
	"net/http"
	"time"

	"github.com/emicklei/structexplorer"
)

func main() {
	log.Println("use service as handler on http://localhost:8080")
	http.ListenAndServe(":8080", structexplorer.NewService("test", time.Now()))
}
