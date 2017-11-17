package main

import (
	"strings"
	"flag"
	"log"
	"net/http"
	"os"
)

// Returned errors
const (
	EOPENLOG = iota
	EWRONGARGS
	EWRONGDIR
)

const logFile = "previewgen.log"

func main() {
	log.Println("start")
	
	// read settings
	var port, dir string
	flag.StringVar(&port, "p", "", "Server port")
	flag.StringVar(&dir, "d", "", "Working directory")
	flag.Parse()
	
	if (len(port) == 0) || (len(dir) == 0) {
		flag.PrintDefaults()
		os.Exit(EWRONGARGS)
	}

	strings.TrimSuffix(dir, "/")
	err := os.MkdirAll(dir, 0770)
	if err != nil {
		log.Println("Cannot create working dir: ", dir, " Error: ", err.Error())
		os.Exit(EOPENLOG)
	}

	// open log file
	str := dir + "/" + logFile
	file, err := os.Create(str)
	if err != nil {
		log.Println("Cannot open log file: ", str, " Error: ", err.Error())
		os.Exit(EOPENLOG)
	}

	log.SetOutput(file)

	// serve http
	handler := &Server{workDir : dir}
	log.Fatal( http.ListenAndServe(":" + port, handler) )  
}

