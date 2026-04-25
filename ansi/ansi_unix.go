//go:build !windows

package ansi

// EnableANSI 在非 Windows 平台上不需要特殊处理
func EnableANSI() {
	// Unix 系统默认支持 ANSI
}
