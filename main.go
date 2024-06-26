package main

import (
	"os"
	"runtime"
	"wujingquan/cvm/commands"
	"wujingquan/cvm/theme"
)

func main() {
	args := os.Args[1:]

	os := runtime.GOOS

	if os != "windows" {
		theme.Error("cvm currently only works on Windows.")
		return
	}

	if len(args) == 0 {
		commands.Help(false)
		return
	}

	switch args[0] {
	case "help":
		commands.Help(false)
	case "list":
		commands.List()
	case "path":
		commands.Path()
	case "install":
		commands.Install(args)
	case "use":
		commands.Use(args[1:])
	default:
		commands.Help(true)
	}
}
