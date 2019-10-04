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
	// Kubernetes username and secrets files and mount path
	userFile string
	secretFile string
	secretMountPath string
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

	flag.StringVar(&userFile, "user", "basic-auth-user", "name of the Kubernetes file containing the username")
	flag.StringVar(&secretFile, "secret", "basic-auth-password", "name of the Kubernetes file containing the password")
	flag.StringVar(&secretMountPath, "mountpath", "/var/secrets", "mount path of the Kubernetes secret used")
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

	if port == 0 || handlerPath == "" || statusPath == "" {
		fmt.Println("Please specify valid handler endpoints and port")
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

	credentialsReader := readBasicAuthFromDisk{
		SecretMountPath:  secretMountPath,
		UserFilename:     userFile,
		PasswordFilename: secretFile,
	}

	credentials, err := credentialsReader.Read()
	if err != nil {
		log.Fatalf("unable to read basic auth credentials: %v", err)
	}

	// listen on all interfaces (":") + port specified
	host := ":" + fmt.Sprintf("%d", port)
	mux := http.NewServeMux()
	mux.Handle("/"+filehandler, decorateWithBasicAuth(status(logfile), credentials))
	mux.Handle("/"+statushandler, decorateWithBasicAuth(status(statusfile), credentials))
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
