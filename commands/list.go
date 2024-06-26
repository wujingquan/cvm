package commands

import (
	"log"
	"os"
	"wujingquan/cvm/theme"

	"github.com/fatih/color"
)

func List() {
	// get users home dir
	homeDir, err := os.UserHomeDir()

	if err != nil {
		log.Fatalln(err)
	}

	// check if .cvmf older exists
	if _, err := os.Stat(homeDir + "/.cvm"); os.IsNotExist(err) {
		theme.Error("No Composer versions installed")
		return
	}

	// check if .cvm/versions folder exists
	if _, err := os.Stat(homeDir + "/.cvm/versions"); os.IsNotExist(err) {
		theme.Error("No Composer versions installed")
		return
	}

	// get all folders in .cvm/versions
	versions, err := os.ReadDir(homeDir + "/.cvm/versions")
	if err != nil {
		log.Fatalln(err)
	}

	theme.Title("Installed Composer versions")

	// print all folders
	for _, version := range versions {
		color.White("    " + version.Name())
	}
}
