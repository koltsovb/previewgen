package main

import (
	"log"	
	"net/http"
	"strings"
)

// Server ...
type Server struct {
	workDir string
}


func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("start request")
	defer log.Println("end request")

	// split resource into parts
	path := strings.Trim(r.URL.Path, "/")
	parts := strings.Split(path, "/")

	
	switch parts[0] {
	case "files": 
		s.routeFiles(w, r)

	default: // return 404
		log.Println("Unknown resource: ", r.URL.Path)
		http.Error(w, "", http.StatusNotFound)
	}

}

func (s *Server) routeFiles(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)

	if r.Method != "POST" {
		log.Println("Unavailable http method: ", r.Method)
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}


	contentType := r.Header.Get("Content-Type")

	switch  {
	case strings.HasPrefix(contentType, "multipart/"):
		s.routeFormFiles(w, r)

	case contentType == "application/json":
		s.routeJSONFiles(w, r)

	default:
		log.Println("Unavailable content-type: ", contentType)
		http.Error(w, "", http.StatusNotFound)
		return
	}
}
