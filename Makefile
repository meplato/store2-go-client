default: client

.PHONY: client deps

client:
	go build github.com/meplato/store2-go-client/v2/cmd/store
