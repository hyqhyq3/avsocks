
linux: 
	mkdir -p dist/linux/amd64
	mkdir -p dist/linux/386
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o dist/linux/amd64/newsocks
	CGO_ENABLED=0 GOARCH=386 GOOS=linux go build -o dist/linux/386/newsocks
	
windows:
	mkdir -p dist/windows/amd64
	mkdir -p dist/windows/386
	CGO_ENABLED=0 GOARCH=amd64 GOOS=windows go build -o dist/windows/amd64/newsocks
	CGO_ENABLED=0 GOARCH=386 GOOS=windows go build -o dist/windows/386/newsocks
	
darwin:
	mkdir -p dist/darwin/amd64
	GOARCH=amd64 GOOS=darwin go build -o dist/windows/amd64/newsocks

all-dist: linux windows darwin
	