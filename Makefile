include $(GOROOT)/src/Make.inc

TARG=goarchive
GOFILES=goarchive.go
GOFMT=gofmt -l -w

CLEANFILES+=./tmp

include $(GOROOT)/src/Make.pkg

format:
	${GOFMT} *.go

all:

gotar:
	make -C $@
