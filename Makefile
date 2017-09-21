build:
	go build main.go

test:
	go test ./commands ./dependency ./helper ./scaffold
