
default: precommit

.PHONY: precommit ensure format generate test check vet errcheck lint gosec osv-scanner trivy addlicense vulncheck run

precommit: ensure format generate test check addlicense
	@echo "ready to commit"

ensure:
	go mod tidy -e
	go mod verify
	rm -rf vendor

format:
	find . -type f -name '*.go' -not -path './vendor/*' -exec gofmt -w "{}" +
	go run -mod=mod github.com/incu6us/goimports-reviser/v3 -project-name github.com/bborbe/ping -format -excludes vendor ./...
	go run -mod=mod github.com/segmentio/golines -w --max-len=100 --shorten-comments -l .
	go run -mod=mod github.com/shoenig/go-modtool -w fmt go.mod

generate:
	rm -rf mocks avro
	mkdir -p mocks
	printf '// Package mocks contains generated mock implementations.\npackage mocks\n' > mocks/mocks.go
	go generate -mod=mod ./...

test:
	go test -mod=mod -p=$${GO_TEST_PARALLEL:-1} -cover -race $(shell go list -mod=mod ./... | grep -v /vendor/)

check: vet errcheck lint

vet:
	go vet -mod=mod $(shell go list -mod=mod ./... | grep -v /vendor/)

errcheck:
	go run -mod=mod github.com/kisielk/errcheck -ignore '(Close|Write|Fprint)' $(shell go list -mod=mod ./... | grep -v /vendor/)

lint:
	go run -mod=mod github.com/golangci/golangci-lint/v2/cmd/golangci-lint run --allow-parallel-runners --timeout 10m ./...

gosec:
	go run -mod=mod github.com/securego/gosec/v2/cmd/gosec -fmt=text $(shell go list -mod=mod ./... | grep -v /vendor/)

osv-scanner:
	go run -mod=mod github.com/google/osv-scanner/v2/cmd/osv-scanner scan --recursive .

trivy:
	trivy fs --exit-code 1 --severity HIGH,CRITICAL .

addlicense:
	go run -mod=mod github.com/google/addlicense -c "Benjamin Borbe" -y $$(date +'%Y') -l bsd $$(find . -name "*.go" -not -path './vendor/*')

vulncheck:
	go run -mod=mod golang.org/x/vuln/cmd/govulncheck $(shell go list -mod=mod ./... | grep -v /vendor/)

run:
	sudo go run main.go 8.8.8.8 8.8.4.4 193.101.111.10 192.168.177.1 192.168.180.1 192.168.178.5
