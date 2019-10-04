package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
)

// Basic Auth uses OpenFaaS implementation logic
type basicAuthCredentials struct {
	User     string
	Password string
}

type readBasicAuthFromDisk struct {
	SecretMountPath string
	UserFilename string
	PasswordFilename string
}

// decorateWithBasicAuth enforces basic auth as a middleware with given credentials
func decorateWithBasicAuth(next http.HandlerFunc, credentials *basicAuthCredentials) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		user, password, ok := r.BasicAuth()
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)

		if !ok || !(credentials.Password == password && user == credentials.User) {

			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("invalid credentials"))
			return
		}

		next.ServeHTTP(w, r)
	}
}

func (r *readBasicAuthFromDisk) Read() (*basicAuthCredentials, error) {
	var credentials *basicAuthCredentials

	if len(r.SecretMountPath) == 0 {
		return nil, fmt.Errorf("invalid SecretMountPath specified for reading secrets")
	}

	userKey := "basic-auth-user"
	if len(r.UserFilename) > 0 {
		userKey = r.UserFilename
	}

	passwordKey := "basic-auth-password"
	if len(r.PasswordFilename) > 0 {
		passwordKey = r.PasswordFilename
	}

	userPath := path.Join(r.SecretMountPath, userKey)
	user, userErr := ioutil.ReadFile(userPath)
	if userErr != nil {
		return nil, fmt.Errorf("unable to load %s", userPath)
	}

	userPassword := path.Join(r.SecretMountPath, passwordKey)
	password, passErr := ioutil.ReadFile(userPassword)
	if passErr != nil {
		return nil, fmt.Errorf("Unable to load %s", userPassword)
	}

	credentials = &basicAuthCredentials{
		User:     strings.TrimSpace(string(user)),
		Password: strings.TrimSpace(string(password)),
	}

	return credentials, nil
}
