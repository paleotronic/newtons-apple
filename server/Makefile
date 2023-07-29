PREFIX=/usr/local/bin

install: ./physics-server
	install ./physics-server $(PREFIX)/physics-server

build: $(wildcard *.go)
	CC="musl-gcc" CGO_LDFLAGS="-static" go build -o "physics-server" -ldflags="-linkmode=external" .
