up:
	docker compose up -d

down:
	docker compose down --rmi all

lint:
	golangci-lint run ./...

test:
	go test -v -race -coverprofile=c.out ./... \
	&& go tool cover -html=c.out \
	&& rm c.out

mock:
	go generate ./...
