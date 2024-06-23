test:
	go test

benchmark:
	go test -bench=. -benchtime=5s -run=^# -benchmem