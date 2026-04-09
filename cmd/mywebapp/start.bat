@echo off
chcp 65001 >nul

REM 切换到脚本所在目录
cd /d "%~dp0"

echo =========================================
echo 俄罗斯中间件监控系统（支持数据源切换）
echo =========================================
echo.

REM 检查Go是否安装
where go >nul 2>nul
if %errorlevel% neq 0 (
    echo 错误: Go未安装
    echo 请从 https://golang.org/dl/ 安装Go
    pause
    exit /b 1
)

echo 编译项目...
go build -o mywebapp.exe .

if %errorlevel% neq 0 (
    echo 编译失败
    pause
    exit /b 1
)

echo 编译成功!
echo.
echo 启动服务...
echo =========================================
echo 管理页面: http://localhost:8083/
echo 设备API: http://localhost:8083/api/devices
echo 数据管理API: http://localhost:8083/api/admin/data-source
echo =========================================
echo.
echo 按 Ctrl+C 停止服务
echo.

REM 启动服务
mywebapp.exe