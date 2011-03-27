include $(GOROOT)/src/Make.inc

TARG=goarchive
GOFILES=goarchive.go
GOFMT=gofmt -l -w

include $(GOROOT)/src/Make.pkg

format:
	${GOFMT} .

all:
