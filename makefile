# Makefile is esaier for me - i am sure others can do this better

pkg_dir = $(GOPATH)/pkg/darwin_amd64

LOGMSG = $(pkg_dir)/github.com/seldonsmule/logmsg.a

all: tesla

clean:
	go clean
	rm -f file.json
	rm -f modelx.json

tesla: file.json
	go build

deps:
	go get github.com/seldonsmule/logmsg
	go get github.com/seldonsmule/restapi
	go get github.com/mattn/go-sqlite3
	go get github.com/denisbrodbeck/machineid
	go get golang.org/x/crypto/ssh/terminal

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
	@echo make deps - gets all the dependancies from github


default: usage
