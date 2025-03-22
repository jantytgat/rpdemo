build_run:
	make snapshot
	./dist/rpdemo_darwin_arm64_v8.0/rpdemo --log-level trace

clean:
	sh scripts/clean.sh

snapshot:
	goreleaser release --snapshot --clean

run:
	./dist/rpdemo_darwin_arm64_v8.0/rpdemo --log-level trace