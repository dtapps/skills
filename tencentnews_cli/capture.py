#!/usr/bin/env python3
"""抓包分析脚本 - 用于捕获 CLI 的网络请求"""

from mitmproxy import http
import json

def request(flow: http.HTTPFlow) -> None:
    if "openapi.inews.qq.com" in flow.request.pretty_host:
        print("\n" + "="*80)
        print(f"[REQUEST] {flow.request.method} {flow.request.path}")
        print(f"[URL] {flow.request.pretty_url}")
        print("[HEADERS]:")
        for name, value in flow.request.headers.items():
            if name.lower() in ['authorization', 'api-key', 'x-api-key']:
                value = value[:20] + "***" if len(value) > 20 else "***"
            print(f"  {name}: {value}")
        
        if flow.request.content:
            print("[BODY]:")
            try:
                body = json.loads(flow.request.content)
                print(json.dumps(body, indent=2, ensure_ascii=False))
            except:
                print(flow.request.content.decode('utf-8', errors='ignore')[:500])
        print("="*80)

def response(flow: http.HTTPFlow) -> None:
    if "openapi.inews.qq.com" in flow.request.pretty_host:
        print(f"\n[RESPONSE] {flow.response.status_code}")
        if flow.response.content:
            print("[BODY]:")
            try:
                body = json.loads(flow.response.content)
                print(json.dumps(body, indent=2, ensure_ascii=False)[:2000])
            except:
                print(flow.response.content.decode('utf-8', errors='ignore')[:500])
        print("="*80 + "\n")
