clean:
	sh scripts/clean.sh

snapshot:
	goreleaser release --snapshot --clean

run_daemon:
	make snapshot
	./dist/rpdemo_darwin_arm64_v8.0/rpdemo --log-level trace