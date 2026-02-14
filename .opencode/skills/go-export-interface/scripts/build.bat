@echo off
cd %~dp0
go build -o go-export-interface.exe export/go-export-interface.go
if %errorlevel% neq 0 (
    echo Goコード生成ツールのビルドに失敗しました。
    exit /b %errorlevel%
)

go build -o dependency-check.exe check/dependency-check.go
if %errorlevel% neq 0 (
    echo 依存関係チェックツールのビルドに失敗しました。
    exit /b %errorlevel%
)
echo ビルドが成功しました。
