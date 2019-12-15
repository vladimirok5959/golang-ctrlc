package ctrlc

import (
	"fmt"
)

func clr(str string) string {
	if !IS_WIN_PLATFORM {
		return fmt.Sprintf("\033[1;31m%s\033[0m", str)
	}
	return str
}

func clg(str string) string {
	if !IS_WIN_PLATFORM {
		return fmt.Sprintf("\033[1;32m%s\033[0m", str)
	}
	return str
}

func cly(str string) string {
	if !IS_WIN_PLATFORM {
		return fmt.Sprintf("\033[1;33m%s\033[0m", str)
	}
	return str
}
