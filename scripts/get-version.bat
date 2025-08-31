@echo off
setlocal enabledelayedexpansion

REM 检查参数
if "%1"=="" (
    echo 用法: %0 ^<major.minor^> [tag]
    echo 例如: %0 4.13
    echo      %0 4.13 tag  ^(同时创建Git标签^)
    exit /b 1
)

set MAJOR_MINOR=%1

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

REM 获取匹配的标签并找到最大构建号
set MAX_BUILD_NUM=0
for /f "tokens=*" %%a in ('git tag -l "!VERSION_PREFIX!-beta*" 2^>nul') do (
    call :extract_build_number "%%a"
)

REM 下一个构建号
set /a NEXT_BUILD_NUM=!MAX_BUILD_NUM! + 1

REM 格式化为4位数字
if !NEXT_BUILD_NUM! lss 10 (
    set BUILD_NUM=000!NEXT_BUILD_NUM!
) else if !NEXT_BUILD_NUM! lss 100 (
    set BUILD_NUM=00!NEXT_BUILD_NUM!
) else if !NEXT_BUILD_NUM! lss 1000 (
    set BUILD_NUM=0!NEXT_BUILD_NUM!
) else (
    set BUILD_NUM=!NEXT_BUILD_NUM!
)
goto create_version

REM 方法2: 使用本地文件管理版本号（备用方案）
:get_version_from_local_file
set VERSION_FILE=.version-counter
set TODAY_KEY=!YEAR!!MONTH_DAY!

if exist "!VERSION_FILE!" (
    REM 读取现有文件
    set /p STORED_DATE=<"!VERSION_FILE!"
    if "!STORED_DATE!"=="!TODAY_KEY!" (
        REM 同一天，读取构建号并递增
        set /p BUILD_NUM_INT=<nul 2>nul
        for /f "skip=1 tokens=*" %%a in (!VERSION_FILE!) do (
            set BUILD_NUM_INT=%%a
            goto found_build_num
        )
        :found_build_num
        if "!BUILD_NUM_INT!"=="" set BUILD_NUM_INT=0
        set /a BUILD_NUM_INT=!BUILD_NUM_INT! + 1
    ) else (
        REM 新的一天，重置为1
        set BUILD_NUM_INT=1
    )
) else (
    REM 文件不存在，从1开始
    set BUILD_NUM_INT=1
)

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

REM 更新本地文件
echo !TODAY_KEY!>"!VERSION_FILE!"
echo !BUILD_NUM_INT!>>"!VERSION_FILE!"

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
goto :eof

REM 函数：从标签中提取构建号
:extract_build_number
set "tag=%~1"
REM 移除前缀，只保留构建号
set "temp=!tag:*beta=!"
REM 转换为数字进行比较
set /a "current_num=!temp!"
if !current_num! gtr !MAX_BUILD_NUM! (
    set MAX_BUILD_NUM=!current_num!
)
goto :eof
