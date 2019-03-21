# About

Minimalistic web server to serve a static files, e.g. a log and status file under separate handlers. Useful for building appliances and other systems where small footprint is preferred.

> **Note:** Right now two handlers are supported (i.e. two files). Open an issue and I'll work on improving this code ;)

See the `Releases` section for pre-compiled builds and source code. Docker artifacts are also available, e.g. `embano1/tinywww:latest`.

## Usage

```bash
Usage of tinywww:
  -handler string
        Path where to register the file http handler (default "/bootstrap")
  -logfile string
        Path to the bootstrap log file (default "/var/log/bootstrap.log")
  -port uint
        Port to listen on (default 8100)
  -status string
        Path where to register the status http handler (default "/status")
  -statusfile string
        Path to the status file (default "/etc/issue")
  -v    Print version information

```

## Run the Binary

```bash
$ ./tinywww -logfile example.file -handler "/log" -statusfile "/etc/issue" -status "/status"
2019/03/21 21:23:06 serving file example.file on ":8100/log"
2019/03/21 21:23:06 serving file /etc/issue on ":8100/status"
$ curl localhost:8100/status
[...]
```

## systemd Example Unit File

An example systemd unit file can be found [here](tinywww.service). Borrowed and modified based on [this](https://medium.com/@benmorel/creating-a-linux-service-with-systemd-611b5c8b91d6) excellent blog post.

> **Note:** Please modify the file as per your needs, e.g. paths, restart behavior, dependencies, etc.