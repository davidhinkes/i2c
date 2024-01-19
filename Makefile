CC = aarch64-linux-gnu-gcc
TARGET = disk

.PHONY: run example all

all: example
run: example
	scp example $(TARGET):. && ssh $(TARGET) "./example"
example:
	CC=$(CC) GOOS=linux GOARCH=arm64 CGO_ENABLED=1 go build ./cmd/example