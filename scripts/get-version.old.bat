@echo off
REM Advanced version management script for Windows
REM Usage: get-version.bat [major.minor]
REM Example: get-version.bat 4.13

setlocal enabledelayedexpansion

REM 默认主版本号
if "%1"=="" (
    set MAJOR_MINOR=4.13
) else (
    set MAJOR_MINOR=%1
)

REM 获取当前日期
for /f "tokens=*" %%a in ('powershell -Command "Get-Date -Format 'yyyy'"') do set YEAR=%%a
for /f "tokens=*" %%a in ('powershell -Command "Get-Date -Format 'MMdd'"') do set MONTH_DAY=%%a

REM 版本前缀格式：4.13.2025.0831
set VERSION_PREFIX=!MAJOR_MINOR!.!YEAR!.!MONTH_DAY!

REM 方法1: 使用Git标签管理版本号
git rev-parse --git-dir >nul 2>&1
if errorlevel 1 (
    REM Git不可用，使用本地文件方法
    goto get_version_from_local_file
)

REM 尝试获取匹配的标签
set LAST_TAG=
for /f "tokens=*" %%a in ('git tag -l "!VERSION_PREFIX!-*" 2^>nul ^| sort') do set LAST_TAG=%%a

if "!LAST_TAG!"=="" (
    REM 没有找到标签，从0001开始
    set BUILD_NUM=0001
    goto create_version
) else (
    REM 从最后一个标签提取构建号
    REM 使用 call 来处理变量替换
    call :extract_build_number "!LAST_TAG!" "!VERSION_PREFIX!-beta"
    set /a BUILD_NUM_INT=!TEMP_TAG! + 1
    
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
    goto create_version
)

REM 方法2: 使用本地文件管理版本号（备用方案）
:get_version_from_local_file
set VERSION_FILE=.version-counter
set TODAY_KEY=!YEAR!!MONTH_DAY!

if exist "!VERSION_FILE!" (
    REM 读取存储的日期和计数
    set /p STORED_DATE=<"!VERSION_FILE!"
    for /f "skip=1" %%a in ('type "!VERSION_FILE!"') do set STORED_COUNT=%%a
    
    if "!STORED_DATE!"=="!TODAY_KEY!" (
        REM 同一天，递增计数
        set /a BUILD_NUM_INT=!STORED_COUNT! + 1
    ) else (
        REM 新的一天，重置为1
        set BUILD_NUM_INT=1
    )
) else (
    REM 首次运行
    set BUILD_NUM_INT=1
)

REM 格式化构建号
if !BUILD_NUM_INT! lss 10 (
    set BUILD_NUM=000!BUILD_NUM_INT!
) else if !BUILD_NUM_INT! lss 100 (
    set BUILD_NUM=00!BUILD_NUM_INT!
) else if !BUILD_NUM_INT! lss 1000 (
    set BUILD_NUM=0!BUILD_NUM_INT!
) else (
    set BUILD_NUM=!BUILD_NUM_INT!
)

REM 更新版本文件
echo !TODAY_KEY!>"!VERSION_FILE!"
echo !BUILD_NUM_INT!>>"!VERSION_FILE!"

set VERSION=!VERSION_PREFIX!-beta!BUILD_NUM!
goto output_version

:create_version
set VERSION=!VERSION_PREFIX!-beta!BUILD_NUM!
goto output_version

:output_version
echo !VERSION!

REM 如果提供了第二个参数 "tag"，则创建Git标签
if "%2"=="tag" (
    git rev-parse --git-dir >nul 2>&1
    if not errorlevel 1 (
        echo 创建Git标签: !VERSION!
        git tag "!VERSION!" 2>nul || echo 标签已存在或无法创建
        
        REM 尝试推送标签到远程
        git ls-remote --exit-code origin >nul 2>&1
        if not errorlevel 1 (
            git push origin "!VERSION!" 2>nul || echo 无法推送标签到远程仓库
        )
    )
)

endlocal
