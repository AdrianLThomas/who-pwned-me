test:
	go test

benchmark:
	go test -bench=. -count=5 -run=^# -benchmem