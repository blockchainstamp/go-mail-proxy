SHELL=PATH='$(PATH)' /bin/sh

PLATFORM := $(shell uname -o)

COMMIT := $(shell $(GIT) rev-parse HEAD)
VERSION ?= $(shell $(GIT) describe --tags ${COMMIT} 2> /dev/null || echo "$(COMMIT)")
BUILD_TIME := $(shell LANG=en_US date +"%F_%T_%z")
ROOT := github.com/blockchainstamp/go-mail-proxy/cmd
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
.PHONY: test

BINDIR=./bin

all: pbs build

build:
	GOOS=$(OS) GOARCH=amd64 $(GOBUILD) -o $(BINDIR)/$(NAME)

mac:
	GOOS=darwin go build  -ldflags="$(LD_FLAGS)"  -o $(BINDIR)/$(NAME).mac $(ROOT)
arm:
	GOOS=linux GOARM=7 GOARCH=arm go build  -ldflags="$(LD_FLAGS)" -o $(BINDIR)/$(NAME).arm
lnx:
	GOOS=linux go build  -ldflags="$(LD_FLAGS)" -o $(NAME).lnx

clean:
	rm $(BINDIR)/$(NAME)
