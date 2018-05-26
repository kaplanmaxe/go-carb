build-resolver:
	go build -o bin/resolver ./cmd/resolver/main.go

run:
	make build-resolver && bin/resolver