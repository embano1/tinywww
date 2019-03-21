# About

Minimalistic web server to serve a static file, e.g. log file. Useful for building appliances and other systems where small footprint is preferred.

See the `Releases` section for pre-compiled builds and source code. Docker artifacts are also available, e.g. `embano1/tinywww:latest`.

## Usage

```bash
Usage of tinywww:
  -file string
        Path to the bootstrap log file (default "/var/log/bootstrap.log")
  -port uint
        Port to listen on (default 8100)
  -v    Print version information
```