build:
	go build -o ${GOPATH}/bin/hws

watch:
	ls **/*.go | entr go build -o ${GOPATH}/bin/hws