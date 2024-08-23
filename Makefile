.DEFAULT_GOAL := help

.PHONY:  help setup install test testacc testcoverage testcleanup terragen docs terraformrc lint ensure-linter ensure-gitleaks ensure-descope ensure-courtney ensure-go
.SILENT: help setup install test testacc testcoverage testcleanup terragen docs terraformrc lint ensure-linter ensure-gitleaks ensure-descope ensure-courtney ensure-go

help: Makefile ## this help message
	grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

setup: install terraformrc ## prepares local environment for running the provider

install: ensure-go ## installs terraform-provider-descope to $GOPATH/bin
	mkdir -p "$$GOPATH/bin"
	go install .
	echo The $$'\e[33m'terraform-provider-descope$$'\e[0m' tool has been installed to $$GOPATH/bin

test: ensure-go ## runs unit tests
	go test -v -timeout 30m $(ARGS) ./...

testacc: ensure-go ## runs acceptance and unit tests
	TF_ACC=1 go test -v -timeout 120m $(ARGS) ./...

testcoverage: ensure-go ensure-courtney ## runs all tests and computes test coverage
	TF_ACC=1 go test -race -timeout 120m -coverpkg=./... -coverprofile=coverage.raw -covermode=atomic $(ARGS) ./...
	cat coverage.raw | grep -v -e "\/tools\/.*\.go\:.*" | grep -v -e ".*\/main\.go\:.*" > coverage.out
	rm -f coverage.raw
	courtney -l coverage.out
	go tool cover -func coverage.out | grep total | awk '{print $$3}'
	go tool cover -html=coverage.out -o coverage.html

testcleanup: ensure-descope ## cleans up local environment after running tests
	descope project list | grep '"name":"Test[^"]\{36,\}"' | sed -e 's/.*"id":"\([^"]*\)".*/\1/' | xargs -I {} descope project delete {} --force

terragen: ensure-go ## runs the terragen tool to generate code and model documentation
	go run tools/terragen/main.go $(ARGS)

docs: ensure-go ## runs tfplugindocs to generate documentation for the registry 
	go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs@v0.19.4 generate -provider-name descope

terraformrc:
	echo 'provider_installation {'                      > ~/.terraformrc
	echo '  dev_overrides {'                            >> ~/.terraformrc
	echo '    "descope/descope" = "'$$GOPATH'/bin"'     >> ~/.terraformrc
	echo '  }'                                          >> ~/.terraformrc
	echo '  direct {}'                                  >> ~/.terraformrc
	echo '}'                                            >> ~/.terraformrc
	echo The $$'\e[33m'.terraformrc$$'\e[0m' file has been created in $$HOME

lint: ensure-linter ensure-gitleaks ## check for linter and gitleaks failures
	golangci-lint --config .github/actions/ci/lint/golangci.yml run
	gitleaks protect --redact -v -c .github/actions/ci/leaks/gitleaks.toml
	gitleaks detect --redact -v -c .github/actions/ci/leaks/gitleaks.toml

ensure-linter: ensure-go
	if ! command -v golangci-lint &> /dev/null; then \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.60.3 ;\
	fi

ensure-gitleaks:
	if ! command -v gitleaks &> /dev/null; then \
		brew install gitleaks ;\
	fi

ensure-descope:
	if ! command -v descope &> /dev/null; then \
	    echo \\nInstall the descope CLI tool with $$'\e[33m'brew install descope$$'\e[0m'\\n ;\
	    false ;\
	fi

ensure-courtney:
	if ! command -v courtney &> /dev/null; then \
	    echo \\nInstall the courtney tool with $$'\e[33m'go install github.com/dave/courtney@master$$'\e[0m'\\n ;\
	    false ;\
	fi

ensure-go:
	if ! command -v go &> /dev/null; then \
	    echo \\nInstall the go compiler from $$'\e[33m'https://go.dev/dl$$'\e[0m'\\n ;\
	    false ;\
	fi
	if [ -z "$$GOPATH" ]; then \
	    echo \\nThe $$'\e[33m'GOPATH$$'\e[0m' environment variable must be defined, see $$'\e[33m'https://go.dev/wiki/GOPATH$$'\e[0m'\\n ;\
	    false ;\
	fi
