run:
	go get -t "github.com/asticode/go-astikit"
	go get -t "github.com/asticode/go-astilectron"
	go get -t "github.com/asticode/go-astilectron-bootstrap"
	cd pkg && go build && pkg.exe

build:
	go get -t "github.com/asticode/go-astikit"
	go get -t "github.com/asticode/go-astilectron"
	go get -t "github.com/asticode/go-astilectron-bootstrap"
	cd pkg && go build