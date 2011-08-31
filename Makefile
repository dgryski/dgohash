include $(GOROOT)/src/Make.inc

TARG=dgohash
GOFILES=\
	superfast.go \
	murmur3.go \
	stringhashes.go \

include $(GOROOT)/src/Make.pkg


