#!/bin/bash
# Advanced version management script
# Usage: ./get-version.sh [major.minor]
# Example: ./get-version.sh 4.13

set -e

# 默认主版本号
MAJOR_MINOR=${1:-"4.13"}

# 获取当前日期
YEAR=$(date +%Y)
MONTH_DAY=$(date +%m%d)

# 版本前缀格式：4.13.2025.0831
VERSION_PREFIX="${MAJOR_MINOR}.${YEAR}.${MONTH_DAY}"

# 方法1: 使用Git标签管理版本号
get_version_from_git_tags() {
    # 获取所有匹配当前日期前缀的标签
    TAGS=$(git tag -l "${VERSION_PREFIX}-*" 2>/dev/null | sort -V)
    
    if [ -z "$TAGS" ]; then
        # 没有找到标签，从0001开始
        BUILD_NUM="0001"
    else
        # 获取最新的标签并提取构建号
        LATEST_TAG=$(echo "$TAGS" | tail -1)
        LAST_BUILD=$(echo "$LATEST_TAG" | sed "s/${VERSION_PREFIX}-beta//")
        
        # 递增构建号
        BUILD_NUM=$(printf "%04d" $((10#$LAST_BUILD + 1)))
    fi
    
    echo "${VERSION_PREFIX}-beta${BUILD_NUM}"
}

# 方法2: 使用远程文件管理版本号（如果有Git远程仓库）
get_version_from_remote() {
    # 尝试获取远程标签
    git fetch --tags 2>/dev/null || true
    get_version_from_git_tags
}

# 方法3: 使用本地文件管理版本号（备用方案）
get_version_from_local_file() {
    VERSION_FILE=".version-counter"
    TODAY_KEY="${YEAR}${MONTH_DAY}"
    
    if [ -f "$VERSION_FILE" ]; then
        # 读取文件内容
        STORED_DATE=$(head -1 "$VERSION_FILE" 2>/dev/null || echo "")
        STORED_COUNT=$(tail -1 "$VERSION_FILE" 2>/dev/null || echo "0")
        
        if [ "$STORED_DATE" = "$TODAY_KEY" ]; then
            # 同一天，递增计数
            BUILD_NUM=$(printf "%04d" $((STORED_COUNT + 1)))
        else
            # 新的一天，重置为0001
            BUILD_NUM="0001"
        fi
    else
        # 首次运行
        BUILD_NUM="0001"
    fi
    
    # 更新版本文件
    echo "$TODAY_KEY" > "$VERSION_FILE"
    echo "${BUILD_NUM#0}" >> "$VERSION_FILE"  # 去掉前导零存储
    
    echo "${VERSION_PREFIX}-beta${BUILD_NUM}"
}

# 主函数：选择最佳的版本获取方法
get_version() {
    if git rev-parse --git-dir > /dev/null 2>&1; then
        # 在Git仓库中
        if git ls-remote --exit-code origin > /dev/null 2>&1; then
            # 有远程仓库，使用远程标签
            get_version_from_remote
        else
            # 仅本地Git，使用本地标签
            get_version_from_git_tags
        fi
    else
        # 不在Git仓库中，使用本地文件
        get_version_from_local_file
    fi
}

# 创建Git标签（如果在Git仓库中）
create_git_tag() {
    VERSION=$1
    if git rev-parse --git-dir > /dev/null 2>&1; then
        echo "创建Git标签: $VERSION"
        git tag "$VERSION" 2>/dev/null || echo "标签已存在或无法创建"
        
        # 尝试推送标签到远程（如果有远程仓库）
        if git ls-remote --exit-code origin > /dev/null 2>&1; then
            git push origin "$VERSION" 2>/dev/null || echo "无法推送标签到远程仓库"
        fi
    fi
}

# 获取版本号
VERSION=$(get_version)
echo "$VERSION"

# 如果提供了第二个参数 "tag"，则创建Git标签
if [ "$2" = "tag" ]; then
    create_git_tag "$VERSION"
fi
