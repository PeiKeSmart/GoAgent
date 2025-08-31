@echo off
chcp 65001 > nul
echo ======================================
echo GoAgent 服务日志查看工具
echo ======================================
echo.

:menu
echo 1. 查看最新20条日志
echo 2. 实时监控日志
echo 3. 查看服务启动日志
echo 4. 查看错误日志
echo 5. 退出
echo.
set /p choice=请选择操作 (1-5): 

if "%choice%"=="1" goto latest
if "%choice%"=="2" goto monitor
if "%choice%"=="3" goto startup
if "%choice%"=="4" goto errors
if "%choice%"=="5" goto exit
echo 无效选择，请重新选择
goto menu

:latest
echo.
echo === 最新20条日志 ===
powershell -Command "Get-Content 'goagent.log' -Encoding UTF8 -Tail 20"
echo.
pause
goto menu

:monitor
echo.
echo === 实时监控日志 (按Ctrl+C退出) ===
powershell -Command "Get-Content 'goagent.log' -Encoding UTF8 -Wait"
goto menu

:startup
echo.
echo === 服务启动日志 ===
powershell -Command "Get-Content 'goagent.log' -Encoding UTF8 | Select-String '日志系统已初始化|星尘代理服务启动|已启动' | Select-Object -Last 10"
echo.
pause
goto menu

:errors
echo.
echo === 错误日志 ===
powershell -Command "Get-Content 'goagent.log' -Encoding UTF8 | Select-String '错误|失败|ERROR|WARN' | Select-Object -Last 10"
echo.
pause
goto menu

:exit
echo 退出日志查看工具
