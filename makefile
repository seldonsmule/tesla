# Makefile is esaier for me - i am sure others can do this better

pkg_dir = $(GOPATH)/pkg/darwin_amd64

LOGMSG = $(pkg_dir)/github.com/seldonsmule/logmsg.a

TESLIB = $(pkg_dir)/tesla.a

all: $(TESLIB) tesla_cmd

clean:
	go clean
	rm -f file.json
	rm -f modelx.json
	rm tesla_cmd

$(TESLIB): file.json modelx.json $(LOGMSG) tesla.go
	go build
	go install

tesla_cmd: example/tesla_cmd.go tesla.go
	go build example/tesla_cmd.go

$(LOGMSG): 
	make deps

deps:
	go get github.com/seldonsmule/logmsg
	go get github.com/seldonsmule/restapi
	go get github.com/mattn/go-sqlite3
	go get github.com/denisbrodbeck/machineid
	go get golang.org/x/crypto/ssh/terminal

rmdeps:
	@rm -f $(LOGMSG)
	@rm -rf $(GOPATH)/src/github.com/seldonsmule/logmsg
	@rm -f $(pkg_dir)/github.com/seldonsmule/restapi.a
	@rm -rf $(GOPATH)/src/github.com/seldonsmule/restapi
	@rm -f $(pkg_dir)/github.com/mattn/go-sqlite3.a
	@rm -rf $(GOPATH)/src/github.com/mattn/go-sqlite3
	@rm -f $(pkg_dir)/github.com/denisbrodbeck/machineid.a
	@rm -rf $(GOPATH)/src/github.com/denisbrodbeck/machineid

getjson: file.json
	@echo Built json file

modelx.json: file.json
	@echo Linking modelx.json file.json
	@ln -s file.json modelx.json

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
