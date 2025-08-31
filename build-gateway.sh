#!/bin/bash
# 边缘网关交叉编译脚本
# 用于编译适合不同硬件架构的版本

echo "GoAgent 边缘网关交叉编译脚本"
echo "================================="

# 常见边缘网关架构配置
declare -A TARGETS=(
    ["amd64"]="linux/amd64"           # 标准 x86_64 工控机
    ["arm64"]="linux/arm64"           # ARM64 设备 (树莓派4+, 工业网关)
    ["armv7"]="linux/arm"             # ARM32 设备 (树莓派3等)
    ["mips64le"]="linux/mips64le"     # 龙芯架构
    ["riscv64"]="linux/riscv64"       # RISC-V 架构
)

# 输出目录
OUTPUT_DIR="./dist"
mkdir -p "$OUTPUT_DIR"

echo "可用的目标架构："
for arch in "${!TARGETS[@]}"; do
    echo "  - $arch (${TARGETS[$arch]})"
done
echo

# 选择要编译的架构
if [ -z "$1" ]; then
    echo "请指定目标架构："
    echo "用法: $0 [amd64|arm64|armv7|mips64le|riscv64|all]"
    exit 1
fi

TARGET_ARCH="$1"

# 编译函数
build_for_arch() {
    local arch="$1"
    local target="${TARGETS[$arch]}"
    local output_file="$OUTPUT_DIR/goagent-$arch"
    
    echo "正在编译 $arch 版本..."
    echo "目标平台: $target"
    
    # 设置交叉编译环境变量
    IFS='/' read -r GOOS GOARCH <<< "$target"
    
    # 根据架构设置特定的编译参数
    case $arch in
        "armv7")
            export GOARM=7  # ARMv7 特定设置
            ;;
        "mips64le")
            export GOMIPS64=hardfloat  # MIPS64 浮点设置
            ;;
    esac
    
    # 执行编译
    env GOOS="$GOOS" GOARCH="$GOARCH" go build \
        -ldflags="-s -w -extldflags '-static'" \
        -tags netgo \
        -o "$output_file" \
        .
    
    if [ $? -eq 0 ]; then
        echo "✅ $arch 编译成功: $output_file"
        # 显示文件信息
        ls -lh "$output_file"
        file "$output_file" 2>/dev/null || echo "文件类型检测不可用"
    else
        echo "❌ $arch 编译失败"
        return 1
    fi
    echo
}

# 执行编译
if [ "$TARGET_ARCH" = "all" ]; then
    echo "编译所有架构版本..."
    for arch in "${!TARGETS[@]}"; do
        build_for_arch "$arch"
    done
elif [ -n "${TARGETS[$TARGET_ARCH]}" ]; then
    build_for_arch "$TARGET_ARCH"
else
    echo "❌ 不支持的架构: $TARGET_ARCH"
    echo "支持的架构: ${!TARGETS[*]}"
    exit 1
fi

echo "编译完成！输出目录: $OUTPUT_DIR"
echo "部署提示："
echo "1. 将对应架构的可执行文件复制到目标设备"
echo "2. 赋予执行权限: chmod +x goagent-[arch]"
echo "3. 以 root 权限安装: sudo ./goagent-[arch] install"
