package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	// where to serve the bootstrap log
	logHandlerPath = "/bootstrap"
	// protect the http server with some sane timeouts
	httpTimeout = 10 * time.Second
	// log file (not path) to serve
	file string
	// port to listen on
	port uint
	// print version information
	v bool
	// release information
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {

	flag.StringVar(&file, "file", "/var/log/bootstrap.log", "Path to the bootstrap log file")
	flag.UintVar(&port, "port", 8100, "Port to listen on")
	flag.BoolVar(&v, "v", false, "Print version information")
	flag.Parse()

	if v {
		fmt.Printf("version: %s-%s-%s", version, commit, date)
		os.Exit(0)
	}

	if file == "" || port == 0 {
		fmt.Println("Please specify a valid file name and port")
		flag.Usage()
		os.Exit(1)
	}

	// check if file exists
	_, err := os.Open(file)
	if err != nil {
		log.Fatalf("error opening file for reading: %v", err)
	}

	// listen on all interfaces (":") + port specified
	host := ":" + fmt.Sprintf("%d", port)
	mux := http.NewServeMux()
	mux.Handle(logHandlerPath, status(file))
	srv := http.Server{
		Addr:         host,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("starting http server on %s (listening on all interfaces)", host)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatalf("error serving HTTP: %v", err)
	}
}

func status(f string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, file)
	}
}
