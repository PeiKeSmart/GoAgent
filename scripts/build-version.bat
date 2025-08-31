@echo off
REM Enhanced build script for GoAgent with dynamic version injection
REM Usage: build-version.bat [target] [version]
REM   target: windows, linux, all (default: windows)
REM   version: version number (default: auto-generated)

setlocal

REM 设置默认参数
if "%1"=="" (
    set TARGET=windows
) else (
    set TARGET=%1
)

REM 获取版本信息
set MAJOR_MINOR=4.13
if "%2"=="" (
    REM 自动生成版本号，调用版本管理脚本
    for /f "tokens=*" %%a in ('call "%~dp0get-version.bat" !MAJOR_MINOR!') do set VERSION=%%a
) else (
    set VERSION=%2
)

REM 获取构建时间（按指定格式：2025-08-31 16:55:00）
for /f "tokens=*" %%a in ('powershell -Command "Get-Date -Format 'yyyy-MM-dd HH:mm:ss'"') do set BUILD_TIME=%%a

REM 获取Git信息（如果可用）
git rev-parse --short HEAD >nul 2>&1
if !errorlevel! equ 0 (
    for /f %%a in ('git rev-parse --short HEAD') do set GIT_COMMIT=%%a
    for /f %%a in ('git rev-parse --abbrev-ref HEAD') do set GIT_BRANCH=%%a
) else (
    set GIT_COMMIT=unknown
    set GIT_BRANCH=unknown
)

REM 获取Go版本
for /f "tokens=3" %%a in ('go version') do set GO_VERSION=%%a

REM 构建ldflags（使用引号处理特殊字符）
set "LDFLAGS=-s -w"
set "LDFLAGS=%LDFLAGS% -X 'main.Version=%VERSION%'"
set "LDFLAGS=%LDFLAGS% -X 'main.BuildTime=%BUILD_TIME%'"
set "LDFLAGS=%LDFLAGS% -X 'main.GitCommit=%GIT_COMMIT%'"
set "LDFLAGS=%LDFLAGS% -X 'main.GitBranch=%GIT_BRANCH%'"
set "LDFLAGS=%LDFLAGS% -X 'main.GoVersion=%GO_VERSION%'"

echo ========================================
echo GoAgent 增强构建脚本
echo ========================================
echo 目标平台: %TARGET%
echo 版本号: %VERSION%
echo 构建时间: %BUILD_TIME%
echo Git提交: %GIT_COMMIT%
echo Git分支: %GIT_BRANCH%
echo Go版本: %GO_VERSION%
echo ========================================

if "%TARGET%"=="windows" (
    echo 正在构建 Windows 版本...
    go build -ldflags="%LDFLAGS%" -o GoAgent.exe .
    if %errorlevel% equ 0 (
        echo ✅ Windows 构建完成: GoAgent.exe
        echo 测试版本信息:
        GoAgent.exe version
    ) else (
        echo ❌ Windows 构建失败！
        exit /b 1
    )
) else if "%TARGET%"=="linux" (
    echo 正在构建 Linux 版本...
    set GOOS=linux
    set GOARCH=amd64
    go build -ldflags="%LDFLAGS%" -o goagent .
    if %errorlevel% equ 0 (
        echo ✅ Linux 构建完成: goagent
    ) else (
        echo ❌ Linux 构建失败！
        exit /b 1
    )
) else if "%TARGET%"=="all" (
    echo 正在构建所有平台版本...
    
    REM Windows
    go build -ldflags="%LDFLAGS%" -o GoAgent.exe .
    if %errorlevel% equ 0 (
        echo ✅ Windows 构建完成: GoAgent.exe
    ) else (
        echo ❌ Windows 构建失败！
        exit /b 1
    )
    
    REM Linux
    set GOOS=linux
    set GOARCH=amd64
    go build -ldflags="%LDFLAGS%" -o goagent .
    if %errorlevel% equ 0 (
        echo ✅ Linux 构建完成: goagent
    ) else (
        echo ❌ Linux 构建失败！
        exit /b 1
    )
) else (
    echo ❌ 未知的目标平台: %TARGET%
    echo 支持的平台: windows, linux, all
    exit /b 1
)

echo ========================================
echo 🎉 构建完成！

REM 自动创建Git标签（如果环境变量 AUTO_TAG=1）
if "%AUTO_TAG%"=="1" (
    echo 正在创建Git标签...
    call "%~dp0get-version.bat" !MAJOR_MINOR! tag >nul
    echo ✅ Git标签已创建: !VERSION!
)

echo ========================================
