.PHONY: build clean deploy

build:
	env GOOS=linux go build -ldflags="-s -w" -o bin/getproducts getproducts/main.go getproducts/utils.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/checkprices checkprices/main.go checkprices/utils.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/updateprices updateprices/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/notifyusers notifyusers/main.go

clean:
	rm -rf ./bin

deploy: clean build
	sls deploy --verbose
