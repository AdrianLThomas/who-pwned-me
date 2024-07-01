test:
	go test ./pkg/**

benchmark:
	go test -bench=. -benchtime=10s -run=^# -benchmem