.PHONY: test benchmark run

test:
	go test ./pkg/**

benchmark:
	go test -bench=. -benchtime=10s -run=^# -benchmem

run:
	go run cmd/main.go test/hibp.txt test/wpm.json

build:
	go build -o who-pwned-me cmd/main.go