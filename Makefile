.PHONY: test benchmark run

test:
	go test ./pkg/**

benchmark:
	go test -bench=. -benchtime=10s -run=^# -benchmem

run-compare:
	go run cmd/main.go compare -hibp-path test/hibp.txt -wpm-path test/wpm.json

run-convert:
	go run cmd/main.go convert -provider bitwarden -path test/bitwarden.json

build:
	go build -o who-pwned-me cmd/main.go