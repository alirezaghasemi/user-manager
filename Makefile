vet:
	go vet ./...

lint:
	golangci-lint -c=.golangci.yml run ./...

test:
	go test -v ./...

cover:
	go test -coverprofile=cover.out ./... && go tool cover -func=cover.out