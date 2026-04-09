@echo off
chcp 65001 >nul

echo =========================================
echo API测试脚本
echo =========================================
echo.

echo 1. 测试管理页面...
curl -s -o nul -w "%%{http_code}" http://localhost:8083/
if %errorlevel% equ 0 (
    echo 成功
) else (
    echo 失败
)

echo.
echo 2. 测试设备API...
curl -s -o nul -w "%%{http_code}" http://localhost:8083/api/devices
if %errorlevel% equ 0 (
    echo 成功
) else (
    echo 失败
)

echo.
echo 3. 测试数据管理API...
curl -s -o nul -w "%%{http_code}" http://localhost:8083/api/admin/data-source
if %errorlevel% equ 0 (
    echo 成功
) else (
    echo 失败
)

echo.
echo 4. 获取当前数据源...
curl -s http://localhost:8083/api/admin/data-source
echo.

echo.
echo =========================================
echo 测试完成
echo =========================================
pause