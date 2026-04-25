#!/bin/sh
# postinstall.sh - 安装后脚本

set -e

echo "✓ dtskills 安装完成"
echo ""
echo "使用方法:"
echo "  dtskills --help          查看帮助"
echo "  dtskills config init     初始化配置"
echo "  dtskills config set      设置配置项"
echo ""
echo "配置文件位置:"
echo "  用户级: ~/.config/dtskills/config.yaml"
echo "  系统级: /etc/dtskills/config.yaml"
echo ""

# 创建软链接 (如果还不存在)
if [ ! -f /usr/local/bin/dtskills ] && [ -f /usr/bin/dtskills ]; then
    ln -sf /usr/bin/dtskills /usr/local/bin/dtskills 2>/dev/null || true
fi

exit 0
