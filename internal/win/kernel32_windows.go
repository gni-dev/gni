package win

import "syscall"

var (
	kernel32            = syscall.NewLazyDLL("kernel32.dll")
	procGetModuleHandle = kernel32.NewProc("GetModuleHandleW")
)

func GetModuleHandle() (syscall.Handle, error) {
	ret, _, err := procGetModuleHandle.Call()
	if ret == 0 {
		return syscall.InvalidHandle, err
	}
	return syscall.Handle(ret), nil
}
