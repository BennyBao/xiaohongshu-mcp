#!/bin/bash

set -e

echo "开始构建 macOS ARM64 版本..."

CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o xiaohongshu-mcp-darwin-arm64 .
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o xiaohongshu-login-darwin-arm64 ./cmd/login

echo "构建完成！"
ls -lh xiaohongshu-*-darwin-arm64
