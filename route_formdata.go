package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)



func (s *Server) routeFormFiles(w http.ResponseWriter, r *http.Request) {

	err := r.ParseMultipartForm(0)
	if err != nil {
		log.Println("routeFormFiles. Error occuered while parse form data: ", err.Error())
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	ids := map[string]uint64{} //answer: "fileName : thumbID"

	headers, _ := r.MultipartForm.File["image_field"]
	for _, h := range headers {
		//read file
		f, err := h.Open() 
		if err != nil {
			log.Println("routeFormFiles. Error occuered while open file: ", err.Error())
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		defer f.Close()

		buf, err := ioutil.ReadAll(f)
		if err != nil {
			log.Println("routeFormFiles. Error occuered while read file: ", err.Error())
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		// create thumb
		id, err := s.saveThumb(buf)
		if err != nil {
			log.Println("routeFormFiles. Cannot save thumb: ", err.Error())
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		ids[h.Filename] = id
	}

	answer, err := json.Marshal(ids)
	if err != nil {
		log.Println("routeFormFiles. Error occured while create json-answer: ", err.Error())
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.Write(answer)
	w.Header().Set("Content-Type", "application/json")
}