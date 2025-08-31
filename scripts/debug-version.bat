@echo off
setlocal enabledelayedexpansion

REM 检查参数
if "%1"=="" (
    echo 用法: %0 ^<major.minor^>
    echo 例如: %0 4.13
    exit /b 1
)

set MAJOR_MINOR=%1

REM 获取当前日期
for /f "tokens=*" %%a in ('powershell -Command "Get-Date -Format 'yyyy'"') do set YEAR=%%a
for /f "tokens=*" %%a in ('powershell -Command "Get-Date -Format 'MMdd'"') do set MONTH_DAY=%%a

REM 版本前缀格式：4.13.2025.0831
set VERSION_PREFIX=!MAJOR_MINOR!.!YEAR!.!MONTH_DAY!

echo MAJOR_MINOR=!MAJOR_MINOR!
echo YEAR=!YEAR!
echo MONTH_DAY=!MONTH_DAY!
echo VERSION_PREFIX=!VERSION_PREFIX!

REM 尝试获取匹配的标签
echo.
echo 查找标签模式: "!VERSION_PREFIX!-*"
git tag -l "!VERSION_PREFIX!-*"

echo.
echo 获取最后一个标签:
set LAST_TAG=
for /f "tokens=*" %%a in ('git tag -l "!VERSION_PREFIX!-*" 2^>nul ^| sort') do (
    echo 找到标签: %%a
    set LAST_TAG=%%a
)

echo LAST_TAG=!LAST_TAG!

if "!LAST_TAG!"=="" (
    echo 没有找到标签，使用 0001
    set BUILD_NUM=0001
) else (
    echo 从标签提取构建号...
    echo 原标签: !LAST_TAG!
    echo 移除前缀 "!VERSION_PREFIX!-beta": 
    set TEMP_TAG=!LAST_TAG:!VERSION_PREFIX!-beta=!
    echo TEMP_TAG=!TEMP_TAG!
    set /a BUILD_NUM_INT=!TEMP_TAG! + 1
    echo BUILD_NUM_INT=!BUILD_NUM_INT!
    
    REM 格式化为4位数字
    if !BUILD_NUM_INT! lss 10 (
        set BUILD_NUM=000!BUILD_NUM_INT!
    ) else if !BUILD_NUM_INT! lss 100 (
        set BUILD_NUM=00!BUILD_NUM_INT!
    ) else if !BUILD_NUM_INT! lss 1000 (
        set BUILD_NUM=0!BUILD_NUM_INT!
    ) else (
        set BUILD_NUM=!BUILD_NUM_INT!
    )
)

echo BUILD_NUM=!BUILD_NUM!
echo 最终版本: !VERSION_PREFIX!-beta!BUILD_NUM!
