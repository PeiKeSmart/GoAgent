#!/bin/bash

echo "GoAgent Linux 权限管理测试脚本"
echo "=================================="

echo
echo "1. 检查当前权限状态"
./goagent check-admin

echo
echo "2. 当前用户信息"
echo "当前用户: $(whoami)"
echo "用户ID: $(id -u)"
echo "sudo是否可用: $(which sudo > /dev/null && echo '是' || echo '否')"

echo
echo "3. 测试权限申请功能（安装服务）"
echo "注意：程序会尝试使用sudo自动重新启动"
read -p "按回车键继续..."
./goagent install

echo
echo "4. 如果上面的自动sudo失败，请手动运行："
echo "sudo ./goagent install"

echo
echo "5. 其他服务管理命令："
echo "sudo ./goagent start     # 启动服务"
echo "sudo ./goagent stop      # 停止服务"
echo "sudo ./goagent uninstall # 卸载服务"

echo
echo "测试完成！"
