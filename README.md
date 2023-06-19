# Introduction
logcat tails an Artifactory request log file, looks for valuable information, parses it and creates a billing log out of it.
Can be used for Artifactory Edge nodes which don't support gathering billing logs.

It uses workers to parse the incoming log lines and a writer to write to the output file. The output file is rotated every hour similar
to the billing logs setup in Artifactory Cloud.

If the log file we are reading from does not exist, logcat will wait for it to be created.

Implementation is inspired by:
https://nesv.github.io/golang/2014/02/25/worker-queues-in-go.html

# Getting Started
Building the binary
```bash
go build -o bin/ ./cmd/*
````

# Executing a test
To run the app:
```bash
PWD=$(pwd)
./bin/logcat -file $PWD/files/artifactory-requests.log -outdir $PWD/files
```

Open another terminal and manually add log entries to the artifactory-request.log:
```bash
echo '2023-01-02T01:02:03.456Z|e227ad976927c6c2|1.2.3.4|user1|HEAD|/api/docker/registry-docker-remote/v2/alpine/curl/manifests/latest|200|-1|1234|567|user-agent123' >> $PWD/files/artifactory-request.log
```
