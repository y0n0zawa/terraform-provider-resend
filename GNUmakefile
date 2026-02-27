default: build

build:
	go build -v ./...

install: build
	go install -v ./...

lint:
	golangci-lint run

fmt:
	golangci-lint fmt

vet:
	go vet ./...

sec:
	go install golang.org/x/vuln/cmd/govulncheck@latest
	govulncheck ./...

generate:
	go generate ./...

docs:
	go install github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs@latest
	tfplugindocs generate

test:
	go test -v -count=1 -parallel=4 ./...

testacc:
	TF_ACC=1 go test -v -count=1 -parallel=4 -timeout 120m ./internal/provider/

check: fmt lint vet sec test

.PHONY: default build install lint fmt vet sec generate docs test testacc check
