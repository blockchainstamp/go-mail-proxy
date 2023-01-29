SHELL=PATH='$(PATH)' /bin/sh

PLATFORM := $(shell uname -o)

COMMIT := $(shell git rev-parse HEAD)
VERSION ?= $(shell git describe --tags ${COMMIT} 2> /dev/null || echo "$(COMMIT)")
BUILD_TIME := $(shell LANG=en_US date +"%F_%T_%z")
ROOT := github.com/blockchainstamp/go-mail-proxy
SRC := $(ROOT)/cmd
LD_FLAGS := -X $(ROOT).Version=$(VERSION) -X $(ROOT).Commit=$(COMMIT) -X $(ROOT).BuildTime=$(BUILD_TIME)

NAME := bmproxy.exe
OS := windows

ifeq ($(PLATFORM), Msys)
    INCLUDE := ${shell echo "$(GOPATH)"|sed -e 's/\\/\//g'}
else ifeq ($(PLATFORM), Cygwin)
    INCLUDE := ${shell echo "$(GOPATH)"|sed -e 's/\\/\//g'}
else
	INCLUDE := $(GOPATH)
	NAME=bmproxy
	OS=linux
endif

# enable second expansion
.SECONDEXPANSION:

.PHONY: all
.PHONY: pbs
.PHONY: test

BINDIR=./bin

all: pbs build

build:
	GOOS=$(OS) GOARCH=amd64 $(GOBUILD) -o $(BINDIR)/$(NAME)

pbs:
	cd utils/ && $(MAKE)

mac:
	GOOS=darwin go build  -o $(BINDIR)/$(NAME).mac   -ldflags="$(LD_FLAGS)"  $(SRC)
arm:
	GOOS=linux GOARM=7 GOARCH=arm go build  -o $(BINDIR)/$(NAME).arm   -ldflags="$(LD_FLAGS)"  $(SRC)
lnx:
	GOOS=linux go build -o $(BINDIR)/$(NAME).lnx   -ldflags="$(LD_FLAGS)"  $(SRC)

clean:
	rm $(BINDIR)/$(NAME).*
