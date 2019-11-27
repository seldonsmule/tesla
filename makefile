# Makefile is esaier for me - i am sure others can do this better

all: tesla

clean:
	go clean
	rm -f file.json
	rm -f modelx.json

tesla: file.json
	go build

getjson: file.json
	@echo Built json file

file.json: 
	go run util/buildjson.go
	ln -s file.json modelx.json

optioncodes.md: 
	git clone https://github.com/timdorr/tesla-api.git
	cp tesla-api/docs/vehicle/optioncodes.md ./
	make cleanjson

usage:
	@echo make all - builds application
	@echo make clean - cleans out all builds
	@echo make getjson - Update codes in json file from github repo

default: usage
