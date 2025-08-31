@echo off
chcp 65001 >nul
setlocal

set LOG_FILE=%~dp0..\goagent.log

echo ========================================
echo          GoAgent 日志查看工具
echo ========================================
echo.

:menu
echo 请选择操作:
echo 1. 查看最新20行日志
echo 2. 查看最新50行日志
echo 3. 实时监控日志
echo 4. 查看日志文件信息
echo 5. 清理旧日志
echo 0. 退出
echo.
set /p choice="请输入选择 (0-5): "

if "%choice%"=="1" goto view20
if "%choice%"=="2" goto view50
if "%choice%"=="3" goto monitor
if "%choice%"=="4" goto info
if "%choice%"=="5" goto cleanup
if "%choice%"=="0" goto exit
echo 无效选择，请重新输入
goto menu

:view20
echo.
echo 最新20行日志:
echo ----------------------------------------
powershell -Command "Get-Content '%LOG_FILE%' -Tail 20 -Encoding UTF8"
echo ----------------------------------------
pause
goto menu

:view50
echo.
echo 最新50行日志:
echo ----------------------------------------
powershell -Command "Get-Content '%LOG_FILE%' -Tail 50 -Encoding UTF8"
echo ----------------------------------------
pause
goto menu

:monitor
echo.
echo 实时监控日志 (按 Ctrl+C 停止):
echo ----------------------------------------
powershell -Command "Get-Content '%LOG_FILE%' -Wait -Encoding UTF8"
goto menu

:info
echo.
echo 日志文件信息:
echo ----------------------------------------
if exist "%LOG_FILE%" (
    echo 文件路径: %LOG_FILE%
    powershell -Command "Get-Item '%LOG_FILE%' | Select-Object Name, Length, LastWriteTime | Format-List"
    powershell -Command "Write-Host '文件行数:' -NoNewline; (Get-Content '%LOG_FILE%' | Measure-Object -Line).Lines"
) else (
    echo 日志文件不存在
)
echo ----------------------------------------
pause
goto menu

:cleanup
echo.
echo 清理旧日志备份文件...
echo ----------------------------------------
cd /d "%~dp0.."
for %%f in (goagent.log.*) do (
    echo 删除: %%f
    del "%%f"
)
echo 清理完成
echo ----------------------------------------
pause
goto menu

:exit
echo 再见！
