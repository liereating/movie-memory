@echo off
setlocal
title Movie DB Launcher
rem One-click start: backend (8080) + frontend (5173), each in its own window.

set "ROOT=%~dp0"

call npm -v >nul 2>&1
if errorlevel 1 (
    echo [ERROR] npm not found. Please install Node.js first.
    pause
    exit /b 1
)

if not exist "%ROOT%backend\movies-server.exe" (
    echo [ERROR] backend\movies-server.exe not found.
    echo         Build it first:
    echo         cd backend
    echo         go build -o movies-server.exe .
    pause
    exit /b 1
)

echo Starting backend :8080 ...
start "backend movies-server :8080" /D "%ROOT%backend" cmd /k "movies-server.exe"

if not exist "%ROOT%frontend\node_modules" (
    echo First run: installing frontend dependencies (npm install)...
    pushd "%ROOT%frontend"
    call npm install
    popd
)

echo Starting frontend :5173 ...
start "frontend vite :5173" /D "%ROOT%frontend" cmd /k "npm run dev"

echo.
echo Both services started. Open this URL in your browser:
echo   http://localhost:5173
echo To stop them run stop.bat, or press Ctrl+C in each window.
pause
