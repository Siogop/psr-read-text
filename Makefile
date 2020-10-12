.PHONY: build deploy

build:
	env GOOS=linux go build -ldflags="-s -w" -o bin/getTextInImage getTextInImage/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/uploadImage uploadImage/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/imageCreated imageCreated/main.go

deploy: 
	rm -rf ./bin
	env GOOS=linux go build -ldflags="-s -w" -o bin/getTextInImage getTextInImage/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/uploadImage uploadImage/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/imageCreated imageCreated/main.go
	sls deploy --verbose
