@echo off
REM Batch file to run lint.ps1 with execution policy bypass
REM Usage: lint.bat [options]

REM Add GOPATH/bin to PATH for this session
for /f "tokens=*" %%i in ('go env GOPATH') do set GOPATH=%%i
set PATH=%PATH%;%GOPATH%\bin

powershell -ExecutionPolicy Bypass -File "%~dp0lint.ps1" %*
