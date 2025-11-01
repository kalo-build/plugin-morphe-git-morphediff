@echo off
REM Morphe Git Diff - Generate semantic schema diffs using git refs
REM Usage: diff-with-git.bat [BASE_REF] [HEAD_REF] [MORPHE_PATH] [OUTPUT]

setlocal EnableDelayedExpansion

set BASE_REF=%1
set HEAD_REF=%2
set MORPHE_PATH=%3
set OUTPUT=%4

if "%BASE_REF%"=="" set BASE_REF=main
if "%HEAD_REF%"=="" set HEAD_REF=HEAD
if "%MORPHE_PATH%"=="" set MORPHE_PATH=morphe
if "%OUTPUT%"=="" set OUTPUT=morphe-diff.yaml

echo.
echo Morphe Git Diff
echo    Comparing: %MORPHE_PATH%
echo    Base: %BASE_REF%
echo    Head: %HEAD_REF%
echo.

REM Validate git repository
git rev-parse --git-dir >nul 2>&1
if errorlevel 1 (
    echo Error: Not a git repository
    exit /b 1
)

REM Validate base ref
git rev-parse %BASE_REF% >nul 2>&1
if errorlevel 1 (
    echo Error: Base ref '%BASE_REF%' not found
    exit /b 1
)

REM Create temp directory
set TEMP_BASE=%TEMP%\morphe-diff-base-%RANDOM%
mkdir "%TEMP_BASE%" 2>nul

REM Extract base from git
echo Extracting base from %BASE_REF%...
git archive %BASE_REF% %MORPHE_PATH% | tar -x -C "%TEMP_BASE%"
if errorlevel 1 (
    echo Error: Failed to extract %MORPHE_PATH% from %BASE_REF%
    rmdir /s /q "%TEMP_BASE%"
    exit /b 1
)

REM Determine head path
if "%HEAD_REF%"=="HEAD" (
    set HEAD_PATH=.\%MORPHE_PATH%
    echo Using working directory for head
) else (
    set TEMP_HEAD=%TEMP%\morphe-diff-head-%RANDOM%
    mkdir "!TEMP_HEAD!" 2>nul
    echo Extracting head from %HEAD_REF%...
    git archive %HEAD_REF% %MORPHE_PATH% | tar -x -C "!TEMP_HEAD!"
    if errorlevel 1 (
        echo Error: Failed to extract %MORPHE_PATH% from %HEAD_REF%
        rmdir /s /q "%TEMP_BASE%"
        rmdir /s /q "!TEMP_HEAD!"
        exit /b 1
    )
    set HEAD_PATH=!TEMP_HEAD!\%MORPHE_PATH%
)

REM Run diff plugin
echo Generating diff...

REM Build config JSON (Windows compatible)
set CONFIG={"baseInputPath":"%TEMP_BASE%\\%MORPHE_PATH%","headInputPath":"%HEAD_PATH%","outputPath":"%OUTPUT%","verbose":false}

REM Execute via kalo
kalo compile --plugin @kalo-build/plugin-morphe-git-morphediff --config "%CONFIG%"
if errorlevel 1 (
    echo Error: Diff generation failed
    rmdir /s /q "%TEMP_BASE%"
    if defined TEMP_HEAD rmdir /s /q "!TEMP_HEAD!"
    exit /b 1
)

REM Cleanup
rmdir /s /q "%TEMP_BASE%"
if defined TEMP_HEAD rmdir /s /q "!TEMP_HEAD!"

echo.
echo Diff generated: %OUTPUT%
echo.

REM Show summary if file exists
if exist "%OUTPUT%" (
    echo Summary:
    findstr /C:"total_changes:" "%OUTPUT%"
    findstr /C:"breaking:" "%OUTPUT%"
    findstr /C:"additive:" "%OUTPUT%"
)

endlocal

