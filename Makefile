fmt:
	go fmt ./pkg/... ./cmd

vet:
	go vet ./pkg/... ./cmd

release: fmt vet
	go build -ldflags "-s -w" -o bin/qks