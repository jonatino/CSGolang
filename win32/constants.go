package win32

const (
	TH32CS_SNAPHEAPLIST = 0x00000001
	TH32CS_SNAPPROCESS  = 0x00000002
	TH32CS_SNAPTHREAD   = 0x00000004
	TH32CS_SNAPMODULE   = 0x00000008
	TH32CS_SNAPALL      = TH32CS_SNAPHEAPLIST | TH32CS_SNAPPROCESS | TH32CS_SNAPTHREAD | TH32CS_SNAPMODULE
)

const (
	MAX_MODULE_NAME32 = 255
	MAX_PATH          = 260
)

const (
	INPUT_MOUSE    = 0
	INPUT_KEYBOARD = 1
	INPUT_HARDWARE = 2
)

const (
	ErrSuccess                = "The operation completed successfully."
	ErrPartialReadWriteMemory = "Only part of a ReadProcessMemory or WriteProcessMemory request was completed."
)
