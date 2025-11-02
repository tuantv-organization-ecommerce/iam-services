@echo off
REM Batch file to run verify-lint.ps1 with execution policy bypass
REM Usage: verify.bat [options]

REM Add GOPATH/bin to PATH for this session
for /f "tokens=*" %%i in ('go env GOPATH') do set GOPATH=%%i
set PATH=%PATH%;%GOPATH%\bin

powershell -ExecutionPolicy Bypass -File "%~dp0verify-lint.ps1" %*
