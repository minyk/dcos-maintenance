PWD := `pwd`

default: build

build: clean linux darwin windows

clean:
	rm -f ./dcos-maintenance*

linux:
	docker run --rm -e "GOOS=linux" -e "GOARCH=amd64" -v $(PWD):/usr/src/github.com/minyk/dcos-maintenance -w /usr/src/github.com/minyk/dcos-maintenance golang:1.11 go build -ldflags="-s -w ${GO_LDFLAGS}" -v -o build/dcos-maintenance-linux

darwin:
	docker run --rm -e "GOOS=darwin" -e "GOARCH=amd64" -v $(PWD):/usr/src/github.com/minyk/dcos-maintenance -w /usr/src/github.com/minyk/dcos-maintenance golang:1.11 go build -ldflags="-s -w ${GO_LDFLAGS}" -v -o build/dcos-maintenance-darwin

windows:
	docker run --rm -e "GOOS=linux" -e "GOARCH=amd64" -v $(PWD):/usr/src/github.com/minyk/dcos-maintenance -w /usr/src/github.com/minyk/dcos-maintenance golang:1.11 go build -ldflags="-s -w ${GO_LDFLAGS}" -v -o build/dcos-maintenance-windows.exe
