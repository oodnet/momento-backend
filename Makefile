all: momento

momento:
	go build ./cmd/momento

clean:
	-rm -f momento

.PHONY: all server
