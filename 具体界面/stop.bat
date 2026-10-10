@echo off
setlocal
title Movie DB Stopper
rem Stop backend (8080) + frontend (5173).

echo Stopping backend (movies-server.exe) ...
"%SystemRoot%\System32\taskkill.exe" /f /im movies-server.exe >nul 2>&1

rem Fallback: kill whatever is listening on 8080.
for /f "tokens=5" %%p in ('%SystemRoot%\System32\netstat.exe -ano ^| %SystemRoot%\System32\findstr.exe ":8080 " ^| %SystemRoot%\System32\findstr.exe "LISTENING"') do (
    "%SystemRoot%\System32\taskkill.exe" /f /pid %%p >nul 2>&1
)

echo Stopping frontend (:5173) ...
for /f "tokens=5" %%p in ('%SystemRoot%\System32\netstat.exe -ano ^| %SystemRoot%\System32\findstr.exe ":5173 " ^| %SystemRoot%\System32\findstr.exe "LISTENING"') do (
    "%SystemRoot%\System32\taskkill.exe" /f /pid %%p >nul 2>&1
)

echo All services stopped.
pause