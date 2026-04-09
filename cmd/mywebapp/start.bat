@echo off
echo ========================================
echo 俄罗斯中间件服务器监控系统
echo ========================================
echo.

echo 编译项目...
go build -o mywebapp.exe

if %ERRORLEVEL% neq 0 (
    echo 编译失败!
    pause
    exit /b 1
)

echo 编译成功!
echo.

echo 启动监控服务...
echo 监听端口: 8083
echo.
echo 访问地址:
echo - 管理界面: http://localhost:8083/
echo - 设备监控: http://localhost:8083/devices.html
echo - API文档: http://localhost:8083/
echo.

mywebapp.exe

if %ERRORLEVEL% neq 0 (
    echo 服务启动失败!
    pause
    exit /b 1
)