@echo off
setlocal enabledelayedexpansion

set HOSTS_FILE=%WINDIR%\System32\drivers\etc\hosts

echo Checking if entry already exists...
findstr /C:"127.0.0.1 minio" "%HOSTS_FILE%" >nul
if %errorlevel% equ 0 (
    echo Host entry already exists in %HOSTS_FILE%
    goto :end
)

echo Adding entry to hosts file (requires administrator privileges)...

:: Check if running as administrator
net session >nul 2>&1
if %errorlevel% neq 0 (
    echo This script requires administrator privileges.
    echo Please right-click and select "Run as administrator"
    goto :end
)

:: Add entry to hosts file
echo 127.0.0.1 minio>> "%HOSTS_FILE%"
if %errorlevel% neq 0 (
    echo Failed to write to hosts file.
) else (
    echo Entry added successfully!
)

:end
echo.
echo Done! To test, try: ping minio
timeout /t 5 