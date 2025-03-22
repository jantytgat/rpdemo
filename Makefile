build_run:
	make snapshot
	make run

clean:
	sh scripts/clean.sh

snapshot:
	goreleaser release --snapshot --clean

run:
	go run cmd/rpdemo/main.go daemon --log-level trace