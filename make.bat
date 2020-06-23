@echo off
SETLOCAL

set APP=Housekeeper
set VERSION=1.3.1
set BINARY-X86=%APP%_%VERSION%.windows.386.exe
set BINARY-X64=%APP%_%VERSION%.windows.amd64.exe
set BINARY-LINUX=%APP%_%VERSION%.linux.amd64

REM Set build number from git commit hash
for /f %%i in ('git rev-parse HEAD') do set BUILD=%%i

set LDFLAGS=-ldflags "-X main.version=%VERSION% -X main.build=%BUILD% -s -w -extldflags '-static'"

goto build

:build
    set GOOS=windows

    echo "=== Building x86 ==="
    set GOARCH=386

    go build -mod=vendor -o %BINARY-X86% %LDFLAGS%

    echo "=== Building x64 ==="
    set GOARCH=amd64

    go build -mod=vendor -o %BINARY-X64% %LDFLAGS%

    echo "=== Building Linux x64 ==="
    set GOOS=linux

    go build -mod=vendor -o %BINARY-LINUX% %LDFLAGS%

    goto :finalise

:finalise
    set GOOS=
    set GOARCH=

    goto :EOF
