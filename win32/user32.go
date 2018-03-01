package win32

import (
	// #include <wtypes.h>
	// #include <winable.h>
	"C"
	"errors"
	"fmt"
	"syscall"
	"unsafe"
)

var (
	moduser32               = syscall.NewLazyDLL("user32.dll")
	procFindWindowW         = moduser32.NewProc("FindWindowW")
	procGetClientRect       = moduser32.NewProc("GetClientRect")
	procGetCursorPos        = moduser32.NewProc("GetCursorPos")
	procGetForegroundWindow = moduser32.NewProc("GetForegroundWindow")
	procGetKeyState         = moduser32.NewProc("GetKeyState")
	procGetWindowRect       = moduser32.NewProc("GetWindowRect")
	procPeekMessage         = moduser32.NewProc("PeekMessageW")
	procSetWindowsHookEx    = moduser32.NewProc("SetWindowsHookExW")
	procUnhookWindowsHookEx = moduser32.NewProc("UnhookWindowsHookEx")
	procSendInput           = moduser32.NewProc("SendInput")
)

func FindWindowW(className, windowName *uint16) HWND {
	ret, _, _ := procFindWindowW.Call(
		uintptr(unsafe.Pointer(className)),
		uintptr(unsafe.Pointer(windowName)))

	return HWND(ret)
}

func GetForegroundWindow() (hwnd syscall.Handle, err error) {
	r0, _, e1 := syscall.Syscall(procGetForegroundWindow.Addr(), 0, 0, 0, 0)
	if e1 != 0 {
		err = error(e1)
		return
	}
	hwnd = syscall.Handle(r0)
	return
}

func GetWindowRect(hwnd HWND) *RECT {
	var rect RECT
	procGetWindowRect.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(&rect)))

	return &rect
}

func GetClientRect(hwnd HWND) *RECT {
	var rect RECT
	ret, _, _ := procGetClientRect.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(&rect)))

	if ret == 0 {
		panic(fmt.Sprintf("GetClientRect(%d) failed", hwnd))
	}

	return &rect
}

func PeekMessage(hwnd HWND, wMsgFilterMin, wMsgFilterMax, wRemoveMsg uint32) (msg MSG, err error) {
	_, _, err = procPeekMessage.Call(
		uintptr(unsafe.Pointer(&msg)),
		uintptr(hwnd),
		uintptr(wMsgFilterMin),
		uintptr(wMsgFilterMax),
		uintptr(wRemoveMsg))

	if err.Error() != ErrSuccess {
		return
	}
	err = nil
	return
}

func GetKeyState(vKey int) (keyState uint16) {
	ret, _, _ := procGetKeyState.Call(uintptr(vKey))
	return uint16(ret)
}

func GetCursorPos() (x, y int, ok bool) {
	pt := POINT{}
	ret, _, _ := procGetCursorPos.Call(uintptr(unsafe.Pointer(&pt)))
	return int(pt.X), int(pt.Y), ret != 0
}

func SetWindowsHookEx(idHook int, lpfn HOOKPROC, hMod HINSTANCE, dwThreadId DWORD) HHOOK {
	ret, _, _ := procSetWindowsHookEx.Call(
		uintptr(idHook),
		uintptr(syscall.NewCallback(lpfn)),
		uintptr(hMod),
		uintptr(dwThreadId),
	)
	return HHOOK(ret)
}

func UnhookWindowsHookEx(hhk HHOOK) bool {
	ret, _, _ := procUnhookWindowsHookEx.Call(
		uintptr(hhk),
	)
	return ret != 0
}

func SendInput(inputs []INPUT) (err error) {
	var validInputs []C.INPUT

	for _, oneInput := range inputs {
		input := C.INPUT{_type: C.DWORD(oneInput.Type)}

		switch oneInput.Type {
		case INPUT_MOUSE:
			(*MouseInput)(unsafe.Pointer(&input)).mi = oneInput.Mi
		case INPUT_KEYBOARD:
			(*KbdInput)(unsafe.Pointer(&input)).ki = oneInput.Ki
		case INPUT_HARDWARE:
			(*HardwareInput)(unsafe.Pointer(&input)).hi = oneInput.Hi
		default:
			err = errors.New("Unknown input type passed: " + fmt.Sprintf("%d", oneInput.Type))
			return
		}

		validInputs = append(validInputs, input)
	}

	_, _, err = procSendInput.Call(
		uintptr(len(validInputs)),
		uintptr(unsafe.Pointer(&validInputs[0])),
		uintptr(unsafe.Sizeof(C.INPUT{})),
	)
	if err.Error() != ErrSuccess {
		return
	}
	err = nil
	return
}
