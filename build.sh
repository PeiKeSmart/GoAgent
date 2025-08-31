#!/bin/bash
# Build script for GoAgent - Unix shell script
# Usage: ./build.sh [target]
#   target: windows, linux, all (default: linux)

TARGET=${1:-linux}

echo "Building GoAgent for $TARGET..."

case $TARGET in
    "windows")
        echo "Building Windows version..."
        GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o GoAgent.exe .
        if [ $? -eq 0 ]; then
            echo "Windows build completed: GoAgent.exe"
        else
            echo "Build failed!"
            exit 1
        fi
        ;;
    "linux")
        echo "Building Linux version..."
        go build -ldflags="-s -w" -o goagent .
        if [ $? -eq 0 ]; then
            echo "Linux build completed: goagent"
            chmod +x goagent
        else
            echo "Build failed!"
            exit 1
        fi
        ;;
    "all")
        echo "Building all versions..."
        
        echo "Building Windows version..."
        GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o GoAgent.exe .
        if [ $? -ne 0 ]; then
            echo "Windows build failed!"
            exit 1
        fi
        
        echo "Building Linux version..."
        GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o goagent .
        if [ $? -ne 0 ]; then
            echo "Linux build failed!"
            exit 1
        fi
        chmod +x goagent
        
        echo "All builds completed successfully!"
        echo "- Windows: GoAgent.exe"
        echo "- Linux: goagent"
        ;;
    *)
        echo "Unknown target: $TARGET"
        echo "Usage: $0 [windows|linux|all]"
        exit 1
        ;;
esac

echo "Build process finished."
