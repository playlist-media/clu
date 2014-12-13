TEST?=./...

default: test

build:
	go build

test:
	go test $(TEST) $(TESTARGS) -timeout=10s -parallel=4

updatedeps:
	go get -u -v ./...

savedeps:
	godep save

clean:
	rm clu

.PHONEY: clean default savedeps testbuild updatedeps
