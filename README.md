# random-file-server

[![Docker Version](https://img.shields.io/docker/v/bunetz/random-file-server?sort=date)](https://hub.docker.com/r/bunetz/random-file-server)
[![Docker Pulls](https://img.shields.io/docker/pulls/bunetz/random-file-server)](https://hub.docker.com/r/bunetz/random-file-server)

### **This code has been almost fully generated by Chat GPT.**
## Description
Small docker container which allows downloading random files with limited speed and custom name.

## Usage
Just run the container and optionally specify flags:
```
  -address string
        Address to listen on (default ":8080")
  -maxSizeMB int
        Maximum size of file to serve in MB (default 100000)
```

Download a test file with:

`curl http://localhost:8080/file?size_mb=<size>&mbs=5&filename=test.bin > filename` or by visiting the same url in a browser.

mbs (to set the MB/s download speed) and filename parameters are optional.
