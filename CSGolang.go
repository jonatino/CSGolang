package main

import (
	"./win32"
	"errors"
	"fmt"
	"unsafe"
)

func main() {
	pid, err := processByName("explorer.exe")

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Your process id is", pid)
}

func processByName(name string) (uint32, error) {
	var entry win32.PROCESSENTRY32
	entry.Size = uint32(unsafe.Sizeof(entry))

	snap := win32.CreateToolhelp32Snapshot(win32.TH32CS_SNAPALL, 0)

	if !win32.Process32First(snap, &entry) {
		return 0, errors.New(fmt.Sprintf("Process %s not found!", name))
	}

	for ok := true; ok; ok = win32.Process32Next(snap, &entry) { //Because why would GO have a fucking do while loop
		processName := GetProcessName(&entry)
		if processName == name {
			return entry.ProcessID, nil
		}
	}

	return 0, errors.New(fmt.Sprintf("Process %s not found!", name))
}

func GetProcessName(entry *win32.PROCESSENTRY32) string {
	return string(entry.SzExeFile[:parseTerminatedString(entry.SzExeFile[:])])
}

func parseTerminatedString(n []byte) int {
	for i := 0; i < len(n); i++ {
		if n[i] == 0 {
			return i
		}
	}
	return len(n)
}
