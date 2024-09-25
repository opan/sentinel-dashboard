dev-run:
	go mod tidy
	go mod vendor
	go run main.go

test:
	go test ./...


dev-ui-run:
	cd ui
	yarn run dev
