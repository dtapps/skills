//go:build windows

package ansi

import (
	"os"
	"syscall"
	"unsafe"
)

// EnableANSI 在 Windows 上启用 ANSI 转义序列支持
func EnableANSI() {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")

	// 设置 UTF-8 代码页以支持 Unicode 字符
	setConsoleCP := kernel32.NewProc("SetConsoleCP")
	setConsoleOutputCP := kernel32.NewProc("SetConsoleOutputCP")
	// 65001 是 UTF-8 代码页
	_, _, _ = setConsoleCP.Call(65001)
	_, _, _ = setConsoleOutputCP.Call(65001)

	// 获取标准输出句柄
	handle := syscall.Handle(os.Stdout.Fd())

	// 获取当前控制台模式
	var mode uint32
	getConsoleMode := kernel32.NewProc("GetConsoleMode")
	setConsoleMode := kernel32.NewProc("SetConsoleMode")

	// #nosec G103 - 使用 unsafe.Pointer 调用 Windows API 是必需的
	ret, _, _ := getConsoleMode.Call(uintptr(handle), uintptr(unsafe.Pointer(&mode)))
	if ret == 0 {
		return // 获取失败，不启用 ANSI
	}

	// 启用 ANSI 转义序列支持 (ENABLE_VIRTUAL_TERMINAL_PROCESSING = 0x0004)
	mode |= 0x0004
	_, _, _ = setConsoleMode.Call(uintptr(handle), uintptr(mode))
}
