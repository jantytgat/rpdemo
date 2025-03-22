build_run:
	make snapshot
	make run

clean:
	sh scripts/clean.sh

snapshot:
	goreleaser release --snapshot --clean

run:
	./dist/rpdemo_darwin_arm64_v8.0/rpdemo daemon --log-level trace