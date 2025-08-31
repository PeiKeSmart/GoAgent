#!/bin/bash
# 检查多个边缘网关设备的 GoAgent 服务状态

# 配置文件路径
DEVICE_CONFIG="devices.conf"

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 检查单个设备状态
check_device() {
    local name="$1"
    local ssh_addr="$2"
    
    echo -e "${BLUE}检查 $name ($ssh_addr):${NC}"
    
    # 测试SSH连接
    if ! ssh -o ConnectTimeout=5 -o BatchMode=yes "$ssh_addr" exit 2>/dev/null; then
        echo -e "  ${RED}❌ SSH连接失败${NC}"
        return 1
    fi
    
    # 获取服务状态信息
    local status_info
    status_info=$(ssh "$ssh_addr" << 'EOF'
        # 检查服务是否存在
        if ! systemctl list-unit-files | grep -q "goagent.service"; then
            echo "SERVICE_NOT_FOUND"
            exit 1
        fi
        
        # 获取服务状态
        if systemctl is-active --quiet goagent; then
            echo "ACTIVE"
        else
            echo "INACTIVE"
        fi
        
        # 获取启用状态
        if systemctl is-enabled --quiet goagent; then
            echo "ENABLED"
        else
            echo "DISABLED"
        fi
        
        # 获取启动时间
        systemctl show goagent --property=ActiveEnterTimestamp --value 2>/dev/null | head -1
        
        # 获取内存使用
        systemctl show goagent --property=MemoryCurrent --value 2>/dev/null | head -1
        
        # 获取PID
        systemctl show goagent --property=MainPID --value 2>/dev/null | head -1
        
        # 获取最近的日志
        journalctl -u goagent --lines=3 --no-pager --quiet 2>/dev/null | tail -1
EOF
    )
    
    if [ $? -ne 0 ]; then
        echo -e "  ${RED}❌ 服务检查失败或服务未安装${NC}"
        return 1
    fi
    
    # 解析状态信息
    local lines=($status_info)
    local active_status="${lines[0]}"
    local enabled_status="${lines[1]}"
    local start_time="${lines[2]}"
    local memory_bytes="${lines[3]}"
    local pid="${lines[4]}"
    local last_log="${lines[5]}"
    
    # 显示服务状态
    if [ "$active_status" = "ACTIVE" ]; then
        echo -e "  ${GREEN}✅ 服务运行中${NC}"
    else
        echo -e "  ${RED}❌ 服务未运行${NC}"
    fi
    
    # 显示启用状态
    if [ "$enabled_status" = "ENABLED" ]; then
        echo -e "  ${GREEN}✅ 开机自启已启用${NC}"
    else
        echo -e "  ${YELLOW}⚠️  开机自启已禁用${NC}"
    fi
    
    # 显示详细信息
    if [ "$active_status" = "ACTIVE" ]; then
        echo -e "  ${BLUE}ℹ️  进程PID: $pid${NC}"
        
        # 转换内存使用
        if [ -n "$memory_bytes" ] && [ "$memory_bytes" != "0" ] && [ "$memory_bytes" != "[not set]" ]; then
            local memory_mb=$((memory_bytes / 1024 / 1024))
            echo -e "  ${BLUE}ℹ️  内存使用: ${memory_mb}MB${NC}"
        fi
        
        # 显示启动时间
        if [ -n "$start_time" ] && [ "$start_time" != "[not set]" ]; then
            echo -e "  ${BLUE}ℹ️  启动时间: $start_time${NC}"
        fi
        
        # 显示最近日志
        if [ -n "$last_log" ]; then
            echo -e "  ${BLUE}ℹ️  最近日志: $last_log${NC}"
        fi
    fi
    
    return 0
}

# 获取服务详细信息
get_service_details() {
    local name="$1"
    local ssh_addr="$2"
    
    echo -e "${BLUE}获取 $name 的详细信息...${NC}"
    
    ssh "$ssh_addr" << 'EOF'
        echo "=== 系统信息 ==="
        uname -a
        echo
        
        echo "=== GoAgent 服务状态 ==="
        systemctl status goagent --no-pager
        echo
        
        echo "=== 最近日志 (最后10行) ==="
        journalctl -u goagent --lines=10 --no-pager
        echo
        
        echo "=== 资源使用情况 ==="
        if [ -n "$(systemctl show goagent --property=MainPID --value)" ]; then
            pid=$(systemctl show goagent --property=MainPID --value)
            if [ "$pid" != "0" ]; then
                ps aux | grep $pid | grep -v grep
            fi
        fi
EOF
}

# 主程序
main() {
    echo -e "${BLUE}GoAgent 服务状态检查工具${NC}"
    echo "=========================="
    echo
    
    # 检查配置文件
    if [ ! -f "$DEVICE_CONFIG" ]; then
        echo -e "${RED}❌ 配置文件 $DEVICE_CONFIG 不存在${NC}"
        echo "请先运行 deploy-batch.sh 创建配置文件"
        exit 1
    fi
    
    local total=0
    local active=0
    local inactive=0
    local failed=0
    
    # 读取配置并检查
    while IFS=, read -r name ssh_addr arch install_path; do
        # 跳过注释行和空行
        [[ "$name" =~ ^[[:space:]]*#.*$ ]] && continue
        [[ -z "$(echo "$name" | tr -d '[:space:]')" ]] && continue
        
        # 清理空格
        name=$(echo "$name" | sed 's/^[[:space:]]*//;s/[[:space:]]*$//')
        ssh_addr=$(echo "$ssh_addr" | sed 's/^[[:space:]]*//;s/[[:space:]]*$//')
        
        total=$((total + 1))
        
        if check_device "$name" "$ssh_addr"; then
            # 通过SSH检查实际服务状态
            if ssh "$ssh_addr" "systemctl is-active --quiet goagent" 2>/dev/null; then
                active=$((active + 1))
            else
                inactive=$((inactive + 1))
            fi
        else
            failed=$((failed + 1))
        fi
        echo
    done < "$DEVICE_CONFIG"
    
    # 显示总结
    echo -e "${BLUE}检查结果总结${NC}"
    echo "=============="
    echo -e "总设备数: $total"
    echo -e "${GREEN}运行中: $active${NC}"
    echo -e "${YELLOW}未运行: $inactive${NC}"
    echo -e "${RED}检查失败: $failed${NC}"
    
    # 建议操作
    if [ $inactive -gt 0 ] || [ $failed -gt 0 ]; then
        echo
        echo -e "${YELLOW}💡 建议操作:${NC}"
        if [ $inactive -gt 0 ]; then
            echo "   - 对于未运行的服务，可以尝试: sudo systemctl start goagent"
        fi
        if [ $failed -gt 0 ]; then
            echo "   - 对于检查失败的设备，请检查网络连接和SSH配置"
        fi
    fi
}

# 显示帮助信息
show_help() {
    echo "GoAgent 服务状态检查工具"
    echo
    echo "用法: $0 [选项]"
    echo
    echo "选项:"
    echo "  -h, --help      显示此帮助信息"
    echo "  -c, --config    指定配置文件 (默认: devices.conf)"
    echo "  -d, --details   显示指定设备的详细信息"
    echo "                  格式: -d '设备名称,SSH地址'"
    echo
    echo "示例:"
    echo "  $0                                    # 检查所有设备"
    echo "  $0 -c my-devices.conf                # 使用自定义配置文件"
    echo "  $0 -d '树莓派4,pi@192.168.1.100'      # 查看详细信息"
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
        -d|--details)
            IFS=, read -r name ssh_addr <<< "$2"
            get_service_details "$name" "$ssh_addr"
            exit 0
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
