// +build windows

package console

import (
	"golang.org/x/sys/windows"
)

func InitConsole() (w, h int, err error) {
	var cbi windows.ConsoleScreenBufferInfo
	w, h = -1, -1
	if err = windows.GetConsoleScreenBufferInfo(windows.Stdout, &cbi); err != nil {
		return
	}
	if err = windows.SetConsoleMode(windows.Stdout, windows.ENABLE_PROCESSED_OUTPUT|windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING); err != nil {
		return
	}
	w, h = int(cbi.MaximumWindowSize.X), int(cbi.MaximumWindowSize.Y)
	return
}
