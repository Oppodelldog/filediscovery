setup: ## Install all the build and lint dependencies
	go get -u golang.org/x/tools/cmd/goimports
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | bash -s v1.27.0
	mkdir .bin || mv bin/golangci-lint .bin/golangci-lint && rm -rf bin

ci-goveralls:
	GO111MODULE=off go get github.com/mattn/goveralls
	goveralls -service drone.io -repotoken oTDITtEKMYs32fahITROxsCIU6z6LHXiy
	
test-with-coverage: ## Run all the tests
	rm -f coverage.tmp && rm -f coverage.txt
	echo 'mode: atomic' > coverage.txt && go list ./... | xargs -n1 -I{} sh -c 'go test -race -covermode=atomic -coverprofile=coverage.tmp {} && tail -n +2 coverage.tmp >> coverage.txt' && rm coverage.tmp

test: ## Run all the tests
	go version
	go env
	go list ./... | xargs -n1 -I{} sh -c 'go test -race {}'

cover: test ## Run all the tests and opens the coverage report
	go tool cover -html=coverage.txt

fmt: ## gofmt and goimports all go files
	find . -name '*.go' -not -wholename './vendor/*' | while read -r file; do gofmt -w -s "$$file"; goimports -w "$$file"; done

lint: ## Run the linters
	.bin/golangci-lint run

ci: setup test-with-coverage lint ## Run all the tests and code checks

drone-ci: ci ci-goveralls ## drone.io build
      
build: test ## build 
	go build

# Self-Documented Makefile see https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help:
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help