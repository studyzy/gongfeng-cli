BINARY_NAME=gongfeng
VERSION=$(shell git describe --tags --always --dirty)
LDFLAGS=-ldflags "-X github.com/studyzy/gongfeng-cli/internal/cmd.Version=$(VERSION)"

.PHONY: all build test lint fmt fmt-check vet golangci-lint coverage install tidy ci release-snapshot release-check clean

all: lint test build

build:
	go build $(LDFLAGS) -o $(BINARY_NAME) ./cmd/gongfeng/

# test 与 CI 一致：开启 race 检测与原子计数覆盖率
test:
	go test -race -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -func=coverage.out

# lint 聚合 CI 中的全部检查项
lint: fmt-check vet golangci-lint

# 检查代码是否已被 gofmt / goimports 格式化（不修改文件），与 CI 行为一致
fmt-check:
	@unformatted=$$(gofmt -l .); \
	if [ -n "$$unformatted" ]; then \
		echo "The following files are not gofmt-formatted:"; \
		echo "$$unformatted"; \
		exit 1; \
	fi
	@command -v goimports >/dev/null 2>&1 || go install golang.org/x/tools/cmd/goimports@latest
	@unformatted=$$(goimports -l .); \
	if [ -n "$$unformatted" ]; then \
		echo "The following files need goimports:"; \
		echo "$$unformatted"; \
		exit 1; \
	fi

vet:
	go vet ./...

golangci-lint:
	@command -v golangci-lint >/dev/null 2>&1 || go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
	golangci-lint run --config=.golangci.yml --timeout=5m ./...

# fmt 用于本地一键格式化代码
fmt:
	gofmt -w .
	@command -v goimports >/dev/null 2>&1 || go install golang.org/x/tools/cmd/goimports@latest
	goimports -w .

# coverage 复用 test 产物生成 HTML 报告
coverage: test
	go tool cover -html=coverage.out -o coverage.html

install:
	go install $(LDFLAGS) ./cmd/gongfeng/

tidy:
	go mod tidy

# ci 目标用于本地完整复现 CI 流程
ci: lint test build

# release-check 校验 .goreleaser.yml 配置合法性
release-check:
	@command -v goreleaser >/dev/null 2>&1 || go install github.com/goreleaser/goreleaser/v2@latest
	goreleaser check

# release-snapshot 在本地生成一个未发布的快照构建（不打 tag、不上传），用于在打 tag 前验证产物
release-snapshot:
	@command -v goreleaser >/dev/null 2>&1 || go install github.com/goreleaser/goreleaser/v2@latest
	goreleaser release --snapshot --clean

clean:
	rm -f $(BINARY_NAME) coverage.out coverage.html
	rm -rf dist/
