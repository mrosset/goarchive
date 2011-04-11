include $(GOROOT)/src/Make.inc
TARG=gotar

GOFILES=gotar.go

GOFMT=gofmt -l -w

include $(GOROOT)/src/Make.cmd

test: format install
	${TARG}
format:
	${GOFMT} .
