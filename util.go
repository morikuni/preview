package preview

import (
	"syscall"
	"unsafe"
)

// WindowSize represents window size.
type WindowSize struct {
	Height uint16
	Width  uint16
	xpixel uint16
	ypixel uint16
}

// WindowSizeError is used when WindowSize() failed.
const WindowSizeError = previewError("failed to load window size")

// GetWindowSize get WindowSize by system call.
func GetWindowSize() (*WindowSize, error) {
	ws := new(WindowSize)
	retCode, _, _ := syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(ws)),
	)

	if retCode != 0 {
		return nil, WindowSizeError
	}

	return ws, nil
}
