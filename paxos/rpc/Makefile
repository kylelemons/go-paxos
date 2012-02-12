PROTOFILES=$(wildcard *.proto)
PBGOFILES=$(patsubst %.proto,%.pb.go,$(PROTOFILES))

all: $(PBGOFILES)

$(PBGOFILES) : $(PROTOFILES)
	@protoc --go_out=. $(PROTOFILES)
