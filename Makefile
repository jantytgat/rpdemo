build_run:
	make snapshot
	make run

clean:
	sh scripts/clean.sh

install_templ:
	go install github.com/a-h/templ/cmd/templ@latest
	go get -tool github.com/a-h/templ/cmd/templ@latest

snapshot:
	goreleaser release --snapshot --clean

run:
	templ generate
	go run cmd/rpdemo/main.go --log-level trace
