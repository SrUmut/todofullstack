bin=./bin/app

run: build
	@$(bin)

build:
	@go build -o $(bin) .