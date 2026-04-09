#!/usr/bin/env python3
"""
API测试脚本
用于验证监控系统的API接口是否正常工作
"""

import requests
import json
import sys

BASE_URL = "http://localhost:8083/api"

def test_api(endpoint, method="GET", data=None):
    """测试API接口"""
    url = f"{BASE_URL}/{endpoint}"

    try:
        if method == "GET":
            response = requests.get(url, timeout=5)
        elif method == "POST":
            response = requests.post(url, json=data, timeout=5)
        elif method == "PUT":
            response = requests.put(url, json=data, timeout=5)
        else:
            print(f"❌ 不支持的HTTP方法: {method}")
            return False

        if response.status_code == 200:
            print(f"✅ {method} {endpoint} - 成功")
            try:
                data = response.json()
                print(f"   响应数据: {json.dumps(data, indent=2, ensure_ascii=False)}")
            except:
                print(f"   响应内容: {response.text}")
            return True
        else:
            print(f"❌ {method} {endpoint} - 失败 (状态码: {response.status_code})")
            print(f"   错误信息: {response.text}")
            return False

    except requests.exceptions.ConnectionError:
        print(f"❌ {method} {endpoint} - 连接失败 (服务未启动?)")
        return False
    except Exception as e:
        print(f"❌ {method} {endpoint} - 异常: {str(e)}")
        return False

def run_all_tests():
    """运行所有测试"""
    print("=" * 60)
    print("俄罗斯中间件监控系统 API 测试")
    print("=" * 60)

    tests = [
        ("devices", "GET"),
        ("devices/status", "GET"),
        ("devices/stats", "GET"),
        ("devices/MIDDLEWARE_001/status", "GET"),
        ("server/status", "GET"),
        ("server/stats", "GET"),
        ("status", "GET"),
    ]

    passed = 0
    total = len(tests)

    for endpoint, method in tests:
        if test_api(endpoint, method):
            passed += 1
        print()

    print("=" * 60)
    print(f"测试结果: {passed}/{total} 通过")

    if passed == total:
        print("✅ 所有API测试通过!")
        return True
    else:
        print("❌ 部分API测试失败")
        return False

if __name__ == "__main__":
    if run_all_tests():
        sys.exit(0)
    else:
        sys.exit(1)