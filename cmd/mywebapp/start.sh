#!/bin/bash

# 俄罗斯中间件监控系统启动脚本

echo "========================================="
echo "俄罗斯中间件监控系统"
echo "========================================="

# 检查Go是否安装
if ! command -v go &> /dev/null; then
    echo "错误: Go未安装"
    echo "请从 https://golang.org/dl/ 安装Go"
    exit 1
fi

# 检查是否在正确的目录
if [ ! -f "main.go" ]; then
    echo "错误: 请在 cmd/mywebapp 目录中运行此脚本"
    exit 1
fi

echo "编译项目..."
go build -o mywebapp .

if [ $? -ne 0 ]; then
    echo "编译失败"
    exit 1
fi

echo "编译成功!"
echo ""
echo "启动服务..."
echo "========================================="
echo "管理页面: http://localhost:8083/"
echo "设备API: http://localhost:8083/api/devices"
echo "数据管理API: http://localhost:8083/api/admin/data-source"
echo "========================================="
echo ""
echo "按 Ctrl+C 停止服务"
echo ""

# 启动服务
./mywebapp