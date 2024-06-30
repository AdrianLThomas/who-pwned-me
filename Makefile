test:
	go test

benchmark:
	go test -bench=. -benchtime=10s -run=^# -benchmem