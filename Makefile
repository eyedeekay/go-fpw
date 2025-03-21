VERSION=0.0.8

fmt:
	find . -name '*.go' -exec gofumpt -w -s {} \;

release: fmt
	gothub release -p -u eyedeekay -r "go-fpw" -t v$(VERSION) -n "lib" -d "tag for release"

build: fmt
	go build -o ffapp/ffapp ./ffapp

run: build
	./ffapp/ffapp