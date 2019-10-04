# About

Minimalistic web server to serve a static files, e.g. a log and status file under separate handlers. Useful for building appliances and other systems where small footprint is preferred. Basic auth is supported (and required to be set up, see Kubernetes section [below](#kubernetes-example-deployment)).

> **Note:** Right now two handlers are supported (i.e. two files). Open an issue and I'll work on improving this code ;)

See the `Releases` section for pre-compiled builds and source code. Docker artifacts are also available, e.g. `embano1/tinywww:latest`.

## Usage

```bash
Usage of tinywww:
  -handler string
        Path where to register the file http handler (default "/bootstrap")
  -logfile string
        Path to the bootstrap log file (default "/var/log/bootstrap.log")
  -mountpath string
        mount path of the Kubernetes secret used (default "/var/secrets")
  -port uint
        Port to listen on (default 8100)
  -secret string
        name of the Kubernetes file containing the password (default "basic-auth-password")
  -status string
        Path where to register the status http handler (default "/status")
  -statusfile string
        Path to the status file (default "/etc/issue")
  -user string
        name of the Kubernetes file containing the username (default "basic-auth-user")
  -v    Print version information
```

## Kubernetes Example Deployment

This application is intended to be run inside Kubernetes (can also be run via systemd, see an older description [here](https://github.com/embano1/tinywww/tree/ed922466d1936cd9fb4e9d91789c92ad374d7ecf)). It expects a secret (configurable via flags) to validate the basic auth credentials. An example deployment file can be found in `tinywww-dep.yaml`.

```bash
# generate a random password
PASSWORD=$(head -c 12 /dev/urandom | shasum| cut -d' ' -f1)

kubectl create secret generic basic-auth \
--from-literal=basic-auth-user=admin \
--from-literal=basic-auth-password="$PASSWORD"
```