# Makefile is esaier for me - i am sure others can do this better
#

pkg_dir = $(GOPATH)/pkg/darwin_arm64

TESLIB = $(pkg_dir)/github.com/seldonsmule/tesla.a

all: $(TESLIB) $(GOPATH)/bin/tesla_admin $(GOPATH)/bin/tesla_cmd $(GOPATH)/bin/tesla_chargelevel

clean:
	go clean
	echo $(TESLIB)
	rm -f file.json
	rm -f modelx.json
	rm -f $(TESLIB)
	rm -f $(GOPATH)/bin/tesla_admin
	rm -f $(GOPATH)/bin/tesla_cmd
	rm -f $(GOPATH)/bin/tesla_chargelevel

$(TESLIB): file.json modelx.json db.go tesla.go login.go
	go build
	go install

$(GOPATH)/bin/tesla_admin: example/tesla_admin.go $(TESLIB)
	go build example/tesla_admin.go
	mv tesla_admin $(GOPATH)/bin

$(GOPATH)/bin/tesla_cmd: example/tesla_cmd.go $(TESLIB)
	go build example/tesla_cmd.go
	mv tesla_cmd $(GOPATH)/bin

$(GOPATH)/bin/tesla_chargelevel: example/tesla_chargelevel.go $(TESLIB)
	go build example/tesla_chargelevel.go
	mv tesla_chargelevel $(GOPATH)/bin

deps:
	go get github.com/seldonsmule/logmsg
	go get github.com/seldonsmule/restapi
	go get github.com/mattn/go-sqlite3
	go get github.com/denisbrodbeck/machineid
	go get golang.org/x/crypto/ssh/terminal

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
