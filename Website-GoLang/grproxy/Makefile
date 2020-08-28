#GOPATH := ${PWD}:${GOPATH}
#export GOPATH

default: build

build:
	go build -ldflags '-s' -tags netgo -v -o ./bin/grproxy ./src/grproxy
