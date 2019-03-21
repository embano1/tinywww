package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
	"unicode"
)

var (
	// where to serve the bootstrap log
	handlerPath string
	// where to serve a status page
	statusPath string
	// protect the http server with some sane timeouts
	httpTimeout = 10 * time.Second
	// log file (not path) to serve
	logfile string
	// status file (not path) to serve
	statusfile string
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

	// TODO: move handlers and file to strings and iterate over to reduce dup code
	flag.StringVar(&handlerPath, "handler", "/bootstrap", "Path where to register the file http handler")
	flag.StringVar(&statusPath, "status", "/status", "Path where to register the status http handler")
	flag.StringVar(&logfile, "logfile", "/var/log/bootstrap.log", "Path to the bootstrap log file")
	flag.StringVar(&statusfile, "statusfile", "/etc/issue", "Path to the status file")
	flag.UintVar(&port, "port", 8100, "Port to listen on")
	flag.BoolVar(&v, "v", false, "Print version information")
	flag.Parse()

	if v {
		fmt.Printf("version: %s-%s-%s", version, commit, date)
		os.Exit(0)
	}

	if logfile == "" || statusfile == "" || port == 0 || handlerPath == "" || statusPath == "" {
		fmt.Println("Please specify a valid file names, handler endpoints and port")
		flag.Usage()
		os.Exit(1)
	}

	// clean the handler string
	filehandler := validateHandler(handlerPath)
	statushandler := validateHandler(statusPath)

	// check if files exists
	_, err := os.Open(logfile)
	if err != nil {
		log.Fatalf("error opening %s for reading: %v", logfile, err)
	}

	_, err = os.Open(statusfile)
	if err != nil {
		log.Fatalf("error opening %s for reading: %v", statusfile, err)
	}

	// listen on all interfaces (":") + port specified
	host := ":" + fmt.Sprintf("%d", port)
	mux := http.NewServeMux()
	mux.Handle("/"+filehandler, status(logfile))
	mux.Handle("/"+statushandler, status(statusfile))
	srv := http.Server{
		Addr:         host,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("serving file %s on \"%s/%s\"", logfile, host, filehandler)
	log.Printf("serving file %s on \"%s/%s\"", statusfile, host, statushandler)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatalf("error serving HTTP: %v", err)
	}
}

func status(f string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, f)
	}
}

func validateHandler(s string) string {
	var handler string
	// remove all non unicode letters and numbers
	handler = strings.TrimFunc(s, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})
	// remove whitespaces
	handler = strings.Replace(handler, " ", "", -1)
	handler = strings.ToLower(handler)
	return handler
}
