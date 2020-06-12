
CC=/usr/local/musl/bin/musl-gcc
GO_LDFLAGS='-linkmode external -extldflags "-static"'

.PHONY: all
all: gossl
gossl: gossl.go build.go init.go
	CC=$(CC) go build --ldflags $(GO_LDFLAGS) $^ 
