package main

import (
	"./win32"
	"fmt"
	"syscall"
	"unsafe"
)

func main() {
	processByName()
}

func processByName() {
	var entry win32.PROCESSENTRY32
	entry.Size = uint32(unsafe.Sizeof(entry))

	snap := win32.CreateToolhelp32Snapshot(win32.TH32CS_SNAPALL, 0)
	fmt.Println(snap)

	for win32.Process32Next(snap, &entry) {
		fmt.Println(entry.ProcessID, syscall.UTF16ToString(entry.SzExeFile[:]), entry.SzExeFile[:])
	}
}
