CGO_ENABLED=0

bin/main: deps $(shell find . -name '*.go')
	CGO_ENABLED=0 go build -o bin/main -a cmd/main.go
	ldd bin/main || true

deps: go.mod
	go mod download

run: bin/main
	./bin/main

clean:
	rm -rf bin vendor

docker-local:
	docker build -f build/package/Dockerfile .

podman-local:
	podman build -f build/package/Dockerfile .

.PHONY: deps run clean docker-local podman-local test
