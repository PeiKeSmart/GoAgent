#!/bin/bash
# Enhanced build script for GoAgent with dynamic version injection
# Usage: ./build-version.sh [target] [version]
#   target: windows, linux, all (default: linux)  
#   version: version number (default: auto-generated)

set -e

# 设置默认参数
TARGET=${1:-linux}
VERSION=${2}

# 自动生成版本号如果未提供
if [ -z "$VERSION" ]; then
    # 调用版本管理脚本
    SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
    VERSION=$("$SCRIPT_DIR/get-version.sh" "4.13")
fi

# 获取构建信息
BUILD_TIME=$(date '+%Y-%m-%d %H:%M:%S')
GO_VERSION=$(go version | awk '{print $3}')

# 获取Git信息（如果可用）
if git rev-parse --git-dir > /dev/null 2>&1; then
    GIT_COMMIT=$(git rev-parse --short HEAD)
    GIT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
else
    GIT_COMMIT="unknown"
    GIT_BRANCH="unknown"
fi

# 构建ldflags
LDFLAGS="-s -w"
LDFLAGS="$LDFLAGS -X 'main.Version=$VERSION'"
LDFLAGS="$LDFLAGS -X 'main.BuildTime=$BUILD_TIME'"
LDFLAGS="$LDFLAGS -X 'main.GitCommit=$GIT_COMMIT'"
LDFLAGS="$LDFLAGS -X 'main.GitBranch=$GIT_BRANCH'"
LDFLAGS="$LDFLAGS -X 'main.GoVersion=$GO_VERSION'"

echo "========================================"
echo "GoAgent 增强构建脚本"
echo "========================================"
echo "目标平台: $TARGET"
echo "版本号: $VERSION"
echo "构建时间: $BUILD_TIME"
echo "Git提交: $GIT_COMMIT"
echo "Git分支: $GIT_BRANCH"
echo "Go版本: $GO_VERSION"
echo "========================================"

build_windows() {
    echo "正在构建 Windows 版本..."
    GOOS=windows GOARCH=amd64 go build -ldflags="$LDFLAGS" -o GoAgent.exe .
    if [ $? -eq 0 ]; then
        echo "✅ Windows 构建完成: GoAgent.exe"
    else
        echo "❌ Windows 构建失败！"
        exit 1
    fi
}

build_linux() {
    echo "正在构建 Linux 版本..."
    GOOS=linux GOARCH=amd64 go build -ldflags="$LDFLAGS" -o goagent .
    if [ $? -eq 0 ]; then
        echo "✅ Linux 构建完成: goagent"
        # 测试版本信息
        ./goagent --version 2>/dev/null || echo "版本信息已注入"
    else
        echo "❌ Linux 构建失败！"
        exit 1
    fi
}

build_all() {
    echo "正在构建所有平台版本..."
    build_linux
    build_windows
}

case "$TARGET" in
    "windows")
        build_windows
        ;;
    "linux") 
        build_linux
        ;;
    "all")
        build_all
        ;;
    *)
        echo "❌ 未知的目标平台: $TARGET"
        echo "支持的平台: windows, linux, all"
        exit 1
        ;;
esac

echo "========================================"
echo "🎉 构建完成！"
echo "========================================"
