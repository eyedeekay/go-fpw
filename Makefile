VERSION=0.0.5

fmt:
	gofmt -w -s *.go

release: fmt
	gothub release -p -u eyedeekay -r "go-fpw" -t v$(VERSION) -n "lib" -d "tag for release"
