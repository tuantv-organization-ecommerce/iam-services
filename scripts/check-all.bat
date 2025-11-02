@echo off
REM Batch file to run check-all.ps1 with execution policy bypass
REM Usage: check-all.bat [options]

REM Add GOPATH/bin to PATH for this session
for /f "tokens=*" %%i in ('go env GOPATH') do set GOPATH=%%i
set PATH=%PATH%;%GOPATH%\bin

powershell -ExecutionPolicy Bypass -File "%~dp0check-all.ps1" %*
