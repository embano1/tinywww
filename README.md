# About

Minimalistic web server to serve a static file, e.g. log file. Useful for building appliances and other systems where small footprint is preferred.

See the `Releases` section for pre-compiled builds and source code. Docker artifacts are also available, e.g. `embano1/tinywww:latest`.

## Usage

```bash
Usage of tinywww:
  -file string
        Path to the bootstrap log file (default "/var/log/bootstrap.log")
  -handler string
        Path where to register the http handler (default "/bootstrap")
  -port uint
        Port to listen on (default 8100)
  -v    Print version information
```

## Example

```bash
$ ./tinywww -file example.file -handler "/test"
2019/03/21 15:32:13 serving file example.file on ":8100/test"
$ curl localhost:8100/test
[...]
```