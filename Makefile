default: client

.PHONY: client deps

client:
	go build github.com/meplato/store2-go-client/cmd/store

deps:
	go get github.com/bgentry/go-netrc/netrc
