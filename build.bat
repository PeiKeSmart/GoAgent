@echo off
REM Build script for GoAgent - Windows batch file
REM Usage: build.bat [target]
REM   target: windows, linux, all (default: windows)

if "%1"=="" (
    set TARGET=windows
) else (
    set TARGET=%1
)

echo Building GoAgent for %TARGET%...

if "%TARGET%"=="windows" (
    echo Building Windows version...
    go build -ldflags="-s -w" -o GoAgent.exe .
    if errorlevel 1 (
        echo Build failed!
        exit /b 1
    )
    echo Windows build completed: GoAgent.exe
) else if "%TARGET%"=="linux" (
    echo Building Linux version...
    set GOOS=linux
    set GOARCH=amd64
    go build -ldflags="-s -w" -o goagent .
    if errorlevel 1 (
        echo Build failed!
        exit /b 1
    )
    echo Linux build completed: goagent
) else if "%TARGET%"=="all" (
    echo Building all versions...
    
    echo Building Windows version...
    go build -ldflags="-s -w" -o GoAgent.exe .
    if errorlevel 1 (
        echo Windows build failed!
        exit /b 1
    )
    
    echo Building Linux version...
    set GOOS=linux
    set GOARCH=amd64
    go build -ldflags="-s -w" -o goagent .
    if errorlevel 1 (
        echo Linux build failed!
        exit /b 1
    )
    
    echo All builds completed successfully!
    echo - Windows: GoAgent.exe
    echo - Linux: goagent
) else (
    echo Unknown target: %TARGET%
    echo Usage: build.bat [windows^|linux^|all]
    exit /b 1
)

echo Build process finished.
