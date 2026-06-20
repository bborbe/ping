
include tools.env

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
	go run github.com/incu6us/goimports-reviser/v3@$(GOIMPORTS_REVISER_VERSION) -project-name github.com/bborbe/ping -format -excludes vendor ./...
	go run github.com/segmentio/golines@$(GOLINES_VERSION) -w --max-len=100 --shorten-comments -l .
	go run github.com/shoenig/go-modtool@$(GO_MODTOOL_VERSION) -w fmt go.mod

generate:
	rm -rf mocks avro
	mkdir -p mocks
	printf '// Package mocks contains generated mock implementations.\npackage mocks\n' > mocks/mocks.go
	go generate -mod=mod ./...

# -race=true catches data races but flakes on some CI runners (rare SIGSEGV
# during gexec.Build in cmd/*-style binary smoke tests). Default off; opt in
# via ENABLE_RACE=true for nightly/manual hardening runs.
TESTFLAGS_RACE = -race=false
ifdef ENABLE_RACE
	TESTFLAGS_RACE = -race=true
endif

test:
	go test -mod=mod -p=$${GO_TEST_PARALLEL:-1} -cover $(TESTFLAGS_RACE) $(shell go list -mod=mod ./... | grep -v /vendor/)

check: vet errcheck lint

vet:
	go vet -mod=mod $(shell go list -mod=mod ./... | grep -v /vendor/)

errcheck:
	go run github.com/kisielk/errcheck@$(ERRCHECK_VERSION) -ignore '(Close|Write|Fprint)' $(shell go list -mod=mod ./... | grep -v /vendor/)

lint:
	go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION) run --allow-parallel-runners --timeout 10m ./...

gosec:
	go run github.com/securego/gosec/v2/cmd/gosec@$(GOSEC_VERSION) -fmt=text $(shell go list -mod=mod ./... | grep -v /vendor/)

osv-scanner:
	go run github.com/google/osv-scanner/v2/cmd/osv-scanner@$(OSV_SCANNER_VERSION) scan --recursive .

trivy:
	trivy fs --exit-code 1 --severity HIGH,CRITICAL .

addlicense:
	go run github.com/google/addlicense@$(ADDLICENSE_VERSION) -c "Benjamin Borbe" -y $$(date +'%Y') -l bsd $$(find . -name "*.go" -not -path './vendor/*')

VULNCHECK_IGNORE ?= GO-2026-4923 GO-2026-4514 GO-2022-0470 GO-2026-4772 GO-2026-4771

# Known-benign govulncheck failure modes we swallow. golang.org/x/tools v0.46.0
# panics on packages containing generic *types.TypeParam during SSA analysis
# (govulncheck v1.3.0+ surface via RuntimeTypes/AllFunctions). We treat that as
# "no findings" because the panic happens AFTER the package scan; any real
# vulnerabilities would have been emitted as JSON on stdout before the panic.
# Any OTHER govulncheck failure (network, bad args, permissions) is surfaced.
vulncheck:
	@PKGS="$(shell go list -mod=mod ./... | grep -v /vendor/)"; \
	IGNORE_JSON=$$(printf '%s\n' $(VULNCHECK_IGNORE) | jq -R . | jq -s .); \
	ERR=$$(mktemp); \
	trap 'rm -f "$$ERR"' EXIT INT TERM; \
	OUT=$$(go run golang.org/x/vuln/cmd/govulncheck@$(GOVULNCHECK_VERSION) -format json $$PKGS 2>$$ERR); \
	RC=$$?; \
	if [ $$RC -ne 0 ] && ! grep -q "ForEachElement called on type containing" "$$ERR"; then \
		echo "govulncheck failed (exit $$RC):" >&2; \
		cat "$$ERR" >&2; \
		exit $$RC; \
	fi; \
	REMAIN=$$(printf '%s' "$$OUT" | jq -rs --argjson ignore "$$IGNORE_JSON" \
		'(map(select(.osv != null)) | map({key: .osv.id, value: (.osv.summary // "")}) | from_entries) as $$sum | \
		 map(select(.finding != null) | .finding) | \
		 map(select(.osv as $$o | $$ignore | index($$o) | not)) | \
		 map("\(.osv)\t\(.trace[-1].module)@\(.trace[-1].version) -> \(.fixed_version)\t\($$sum[.osv] // "")") | \
		 unique | .[]'); \
	if [ -n "$$REMAIN" ]; then \
		echo "Unexpected vulnerabilities (ignored: $(VULNCHECK_IGNORE)):"; \
		printf '%s\n' "$$REMAIN" | column -t -s "$$(printf '\t')"; \
		exit 1; \
	else \
		echo "No unignored vulnerabilities found"; \
	fi

run:
	sudo go run main.go 8.8.8.8 8.8.4.4 193.101.111.10 192.168.177.1 192.168.180.1 192.168.178.5
