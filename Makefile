default:
	make build
	make run

build:
	go build -o bin/ioio

run:
	bin/ioio

.PHONY: default,build,run
