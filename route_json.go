package main

import (
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"log"
	"encoding/json"
	"net/http"
)

type filesInfo struct {
	Type string 
	Files json.RawMessage
}

type fileList struct {
	Name string
	MimeType string
	Content []byte
}

type urlList struct {
	URL string
}

func (s *Server) routeJSONFiles(w http.ResponseWriter, r *http.Request) {
	log.Println("")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Error occuered while read request body: ", err.Error())
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	info := filesInfo{}
	err = json.Unmarshal(body, &info)
	if err != nil {
		log.Println("routeJSONFiles. Error occuered while unmarshal json: ", err.Error())
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	switch info.Type {
	case "files":
		s.processFiles(w, r, info.Files)
	case "urls":
		s.processUrls(w, r, info.Files)
	default:
		log.Println("Error: unknown 'type' of query: ", err.Error())
		http.Error(w, "", http.StatusNotFound)
		return
	}
}

func(s *Server) processFiles(w http.ResponseWriter, r *http.Request, data json.RawMessage) {
	var fl []fileList
	err := json.Unmarshal(data, &fl)
	if err != nil {
		log.Println("processFiles. Error occuered while unmarshal json: ", err.Error())
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	ids := map[string]uint64{} //answer: "fileName : thumbID"

	// create thumb for all files
	for _, item := range fl {

		//check mime type
		if (item.MimeType != "image/png") && (item.MimeType != "image/jpeg") {
			log.Println("processFiles. Unsupported mime type: ", err.Error())
			http.Error(w, "", http.StatusUnsupportedMediaType)
			return
		}
		
		id, err := s.saveThumb(item.Content)
		if err != nil {
			log.Println("processFiles. Cannot save thumb: ", err.Error())
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		ids[item.Name] = id
	}

	// if OK, return ids of thumbs
	answer, err := json.Marshal(ids)
	if err != nil {
		log.Println("processFiles. Error occured while create json-answer: ", err.Error())
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.Write(answer)
	w.Header().Set("Content-Type", "application/json")
}

func(s *Server) processUrls(w http.ResponseWriter, r *http.Request, data json.RawMessage) {
	var ul []urlList
	err := json.Unmarshal(data, &ul)
	if err != nil {
		log.Println("processFiles. Error occuered while unmarshal json: ", err.Error())
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	ids := map[string]uint64{} //answer: "fileURL : thumbID"

	// download all url and crete thumb
	for _, item := range ul {

		resp, err := http.Get(item.URL)
		if err != nil {
			log.Println("processUrls. Error occuered while download url: ", item.URL, " err: : ", err.Error())
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Println("processUrls. Error: ", resp.StatusCode, " for url: ", item.URL)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		//check mime type
		mimeType := resp.Header.Get("Content-type")
		if (mimeType != "image/png") && (mimeType != "image/jpeg") {
			log.Println("processUrls. Unsupported mime type")
			http.Error(w, "", http.StatusUnsupportedMediaType)
			return
		}

		buf, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("processUrls. Error occuered while read response body: ", item.URL, " err: ", err.Error())
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		id, err := s.saveThumb(buf)
		if err != nil {
			log.Println("processUrls. Cannot save thumb: ", err.Error())
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		ids[item.URL] = id
	}

	// if OK, return ids of thumbs
	answer, err := json.Marshal(ids)
	if err != nil {
		log.Println("processUrls. Error occured while create json-answer: ", err.Error())
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.Write(answer)
	w.Header().Set("Content-Type", "application/json")
}

