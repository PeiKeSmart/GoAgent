@echo off
setlocal enabledelayedexpansion

echo ==========================================
echo 版本管理系统测试
echo ==========================================

echo.
echo 1. 测试当前版本获取:
call scripts\get-version.bat 4.13

echo.
echo 2. 测试创建标签:
call scripts\get-version.bat 4.13 tag

echo.
echo 3. 测试下一个版本号:
call scripts\get-version.bat 4.13

echo.
echo 4. 查看所有今天的标签:
git tag -l "4.13.2025.0831-*"

echo.
echo 5. 测试不同主版本号:
call scripts\get-version.bat 5.0

echo.
echo 6. 测试自动构建+标签:
set AUTO_TAG=1
call scripts\build-version.bat windows

echo.
echo 7. 最终标签列表:
git tag -l "*2025.0831-*"

echo.
echo ==========================================
echo 测试完成！
echo ==========================================
