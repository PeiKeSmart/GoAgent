#!/bin/bash
# 批量部署 GoAgent 到多个边缘网关设备

# 配置文件路径
DEVICE_CONFIG="devices.conf"

# 检查依赖
check_dependencies() {
    for cmd in scp ssh; do
        if ! command -v $cmd &> /dev/null; then
            echo "❌ 缺少依赖: $cmd"
            exit 1
        fi
    done
}

# 创建示例配置文件
create_sample_config() {
    if [ ! -f "$DEVICE_CONFIG" ]; then
        echo "创建示例配置文件: $DEVICE_CONFIG"
        cat > "$DEVICE_CONFIG" << 'EOF'
# GoAgent 批量部署配置文件
# 格式：设备名称,SSH地址,架构,安装路径
# 支持的架构: amd64, arm64, armv7, mips64le, riscv64

# 示例配置 (请根据实际情况修改)
树莓派4,pi@192.168.1.100,arm64,/home/pi/goagent
工控机1,admin@192.168.1.200,amd64,/opt/goagent
路由器1,root@192.168.1.1,armv7,/usr/bin/goagent
网关2,admin@192.168.1.201,arm64,/usr/local/bin/goagent
EOF
        echo "请编辑 $DEVICE_CONFIG 文件，配置你的设备信息"
        exit 0
    fi
}

# 部署到单个设备
deploy_device() {
    local name="$1"
    local ssh_addr="$2"
    local arch="$3"
    local install_path="$4"
    
    echo "=== 部署到 $name ($ssh_addr) ==="
    
    # 检查编译文件是否存在
    if [ ! -f "dist/goagent-$arch" ]; then
        echo "编译 $arch 版本..."
        if ! ./build-gateway.sh "$arch"; then
            echo "❌ 编译失败，跳过 $name"
            return 1
        fi
    fi
    
    # 测试SSH连接
    if ! ssh -o ConnectTimeout=5 -o BatchMode=yes "$ssh_addr" exit 2>/dev/null; then
        echo "❌ SSH连接失败，跳过 $name"
        return 1
    fi
    
    # 复制文件
    echo "复制文件到设备..."
    if scp -o ConnectTimeout=10 "dist/goagent-$arch" "$ssh_addr:/tmp/goagent" >/dev/null 2>&1; then
        echo "✅ 文件复制成功"
    else
        echo "❌ 文件复制失败，跳过 $name"
        return 1
    fi
    
    # 远程安装
    echo "远程安装服务..."
    ssh "$ssh_addr" << EOF
        # 移动文件到目标位置
        if [ -w "\$(dirname $install_path)" ]; then
            mv /tmp/goagent $install_path
        else
            sudo mv /tmp/goagent $install_path
        fi
        
        # 设置权限
        if [ -O "$install_path" ]; then
            chmod +x $install_path
        else
            sudo chmod +x $install_path
        fi
        
        # 停止可能存在的旧服务
        sudo systemctl stop goagent 2>/dev/null || true
        sudo $install_path uninstall 2>/dev/null || true
        
        # 安装新服务
        sudo $install_path install
        sudo systemctl start goagent
        sudo systemctl enable goagent
        
        # 检查服务状态
        if systemctl is-active --quiet goagent; then
            echo "✅ 服务安装并启动成功"
        else
            echo "❌ 服务启动失败"
            exit 1
        fi
EOF
    
    if [ $? -eq 0 ]; then
        echo "✅ $name 部署成功"
        return 0
    else
        echo "❌ $name 部署失败"
        return 1
    fi
}

# 主程序
main() {
    echo "GoAgent 批量部署工具"
    echo "==================="
    
    # 检查依赖
    check_dependencies
    
    # 创建示例配置文件
    create_sample_config
    
    # 检查配置文件
    if [ ! -f "$DEVICE_CONFIG" ]; then
        echo "❌ 配置文件 $DEVICE_CONFIG 不存在"
        exit 1
    fi
    
    # 统计信息
    local total=0
    local success=0
    local failed=0
    
    # 读取配置并部署
    while IFS=, read -r name ssh_addr arch install_path; do
        # 跳过注释行和空行
        [[ "$name" =~ ^[[:space:]]*#.*$ ]] && continue
        [[ -z "$(echo "$name" | tr -d '[:space:]')" ]] && continue
        
        # 清理空格
        name=$(echo "$name" | sed 's/^[[:space:]]*//;s/[[:space:]]*$//')
        ssh_addr=$(echo "$ssh_addr" | sed 's/^[[:space:]]*//;s/[[:space:]]*$//')
        arch=$(echo "$arch" | sed 's/^[[:space:]]*//;s/[[:space:]]*$//')
        install_path=$(echo "$install_path" | sed 's/^[[:space:]]*//;s/[[:space:]]*$//')
        
        total=$((total + 1))
        
        if deploy_device "$name" "$ssh_addr" "$arch" "$install_path"; then
            success=$((success + 1))
        else
            failed=$((failed + 1))
        fi
        echo
    done < "$DEVICE_CONFIG"
    
    # 显示总结
    echo "批量部署完成！"
    echo "=============="
    echo "总设备数: $total"
    echo "成功部署: $success"
    echo "部署失败: $failed"
    
    if [ $failed -gt 0 ]; then
        echo
        echo "💡 失败原因可能包括："
        echo "   - SSH连接失败 (检查网络和认证)"
        echo "   - 权限不足 (确保有sudo权限)"
        echo "   - 架构不匹配 (检查设备架构)"
        echo "   - 存储空间不足"
    fi
}

# 显示帮助信息
show_help() {
    echo "GoAgent 批量部署工具"
    echo
    echo "用法: $0 [选项]"
    echo
    echo "选项:"
    echo "  -h, --help     显示此帮助信息"
    echo "  -c, --config   指定配置文件 (默认: devices.conf)"
    echo
    echo "配置文件格式:"
    echo "  设备名称,SSH地址,架构,安装路径"
    echo
    echo "示例:"
    echo "  树莓派4,pi@192.168.1.100,arm64,/home/pi/goagent"
    echo "  工控机,admin@192.168.1.200,amd64,/opt/goagent"
}

# 处理命令行参数
while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            show_help
            exit 0
            ;;
        -c|--config)
            DEVICE_CONFIG="$2"
            shift 2
            ;;
        *)
            echo "未知选项: $1"
            show_help
            exit 1
            ;;
    esac
done

# 运行主程序
main
