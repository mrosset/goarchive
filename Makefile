include $(GOROOT)/src/Make.inc

TARG=goarchive
GOFILES=goarchive.go
GOFMT=gofmt -l -w

CLEANFILES+=./tmp ./test

include $(GOROOT)/src/Make.pkg

format:
	${GOFMT} *.go

all:
