#!/bin/sh
# preremove.sh - 卸载前脚本

set -e

echo "正在卸载 dtskills..."

# 清理软链接
if [ -L /usr/local/bin/dtskills ]; then
  rm -f /usr/local/bin/dtskills 2>/dev/null || true
fi

exit 0
