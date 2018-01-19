package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"time"
	"net/http/pprof"
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

	//------------------------------- start http-server
	//subscribe to Signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	// serve http
	m := http.NewServeMux()

	handler := &Server{workDir: dir}
	m.HandleFunc("/", handler.ServeHTTP)

	// for profiler
	runtime.SetBlockProfileRate(1)
	m.HandleFunc("/debug/pprof/", pprof.Index)
	m.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	m.HandleFunc("/debug/pprof/profile", pprof.Profile)
	m.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	m.HandleFunc("/debug/pprof/trace", pprof.Trace)

	server := &http.Server{Addr: ":" + port, Handler: m}

 
	go func() { 
		err := server.ListenAndServe()
		log.Println("Server exit with err: ", err) 
	}()

	//wait Signal
	<-stop
	log.Println("Start shutdown ...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	server.Shutdown(ctx)
	defer cancel()
	log.Println("Server was stopped")
	
}

