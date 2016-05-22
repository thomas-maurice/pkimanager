all: deps format bin

deps:
	go get -d

format:
	for directory in . pkiutils cmd config; do /bin/sh -c "cd $$directory && go fmt"; done;

bin: format
	go build -v

install: format
	go install
