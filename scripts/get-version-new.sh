#!/bin/bash
# Advanced version management script
# Usage: ./get-version.sh [major.minor] [tag]
# Example: ./get-version.sh 4.13
#          ./get-version.sh 4.13 tag

set -e

# 检查参数
if [ $# -eq 0 ]; then
    echo "用法: $0 <major.minor> [tag]"
    echo "例如: $0 4.13"
    echo "      $0 4.13 tag  (同时创建Git标签)"
    exit 1
fi

MAJOR_MINOR="$1"

# 获取当前日期
YEAR=$(date +%Y)
MONTH_DAY=$(date +%m%d)

# 版本前缀格式：4.13.2025.0831
VERSION_PREFIX="${MAJOR_MINOR}.${YEAR}.${MONTH_DAY}"

# 方法1: 使用Git标签管理版本号
get_version_from_git_tags() {
    # 检查是否在Git仓库中
    if ! git rev-parse --git-dir >/dev/null 2>&1; then
        echo "不在Git仓库中，使用本地文件方法"
        get_version_from_local_file
        return
    fi
    
    # 获取匹配的标签并找到最大构建号
    MAX_BUILD_NUM=0
    
    # 使用git tag命令查找匹配的标签
    while IFS= read -r tag; do
        if [ -n "$tag" ]; then
            # 提取构建号 (移除版本前缀和-beta)
            build_num=$(echo "$tag" | sed "s/${VERSION_PREFIX}-beta//")
            
            # 确保是数字
            if [[ "$build_num" =~ ^[0-9]+$ ]]; then
                if [ "$build_num" -gt "$MAX_BUILD_NUM" ]; then
                    MAX_BUILD_NUM="$build_num"
                fi
            fi
        fi
    done < <(git tag -l "${VERSION_PREFIX}-beta*" 2>/dev/null)
    
    # 下一个构建号
    NEXT_BUILD_NUM=$((MAX_BUILD_NUM + 1))
    
    # 格式化为4位数字
    BUILD_NUM=$(printf "%04d" "$NEXT_BUILD_NUM")
}

# 方法2: 使用本地文件管理版本号（备用方案）
get_version_from_local_file() {
    VERSION_FILE=".version-counter"
    TODAY_KEY="${YEAR}${MONTH_DAY}"
    
    if [ -f "$VERSION_FILE" ]; then
        # 读取现有文件
        STORED_DATE=$(head -1 "$VERSION_FILE" 2>/dev/null || echo "")
        if [ "$STORED_DATE" = "$TODAY_KEY" ]; then
            # 同一天，读取构建号并递增
            BUILD_NUM_INT=$(tail -1 "$VERSION_FILE" 2>/dev/null || echo "0")
            BUILD_NUM_INT=$((BUILD_NUM_INT + 1))
        else
            # 新的一天，重置为1
            BUILD_NUM_INT=1
        fi
    else
        # 文件不存在，从1开始
        BUILD_NUM_INT=1
    fi
    
    # 格式化为4位数字
    BUILD_NUM=$(printf "%04d" "$BUILD_NUM_INT")
    
    # 更新本地文件
    echo "$TODAY_KEY" > "$VERSION_FILE"
    echo "$BUILD_NUM_INT" >> "$VERSION_FILE"
}

# 执行版本获取
get_version_from_git_tags

# 构建最终版本号
VERSION="${VERSION_PREFIX}-beta${BUILD_NUM}"

# 输出版本号
echo "$VERSION"

# 如果提供了第二个参数 "tag"，则创建Git标签
if [ "$2" = "tag" ]; then
    if git rev-parse --git-dir >/dev/null 2>&1; then
        echo "创建Git标签: $VERSION"
        if git tag "$VERSION" 2>/dev/null; then
            # 尝试推送标签到远程
            if git ls-remote --exit-code origin >/dev/null 2>&1; then
                git push origin "$VERSION" 2>/dev/null || echo "无法推送标签到远程仓库"
            fi
        else
            echo "标签已存在或无法创建"
        fi
    fi
fi
