CGO_ENABLED=0

bin/main: $(shell find . -name '*.go')
	CGO_ENABLED=0 go build -o bin/main -a cmd/main.go
	ldd bin/main || true

deps: go.mod
	go mod download

run: bin/main
	PROMETHEUS_URI=http://192.168.6.10:9090 \
	WRITE_TO_FILE_PATH=/tmp/out.png \
		./bin/main

clean:
	rm -rf bin vendor

docker-local:
	docker build -f build/package/Dockerfile .

podman-local:
	podman build -f build/package/Dockerfile .

.PHONY: deps run clean docker-local podman-local test
