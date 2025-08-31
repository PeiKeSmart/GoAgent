@echo off
echo GoAgent 管理员权限测试脚本
echo ================================

echo.
echo 1. 检查当前权限状态
GoAgent.exe check-admin

echo.
echo 2. 测试权限申请功能（安装服务）
echo 注意：这将弹出UAC对话框
pause
GoAgent.exe install

echo.
echo 3. 启动服务
GoAgent.exe start

echo.
echo 4. 检查服务状态
sc query GoAgent

echo.
echo 5. 停止服务
GoAgent.exe stop

echo.
echo 6. 卸载服务
GoAgent.exe uninstall

echo.
echo 测试完成！
pause
