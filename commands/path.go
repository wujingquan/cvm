package commands

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"wujingquan/cvm/theme"
)

func Path() {
	theme.Title("cvm: Composer Version Manager")

	// get home dir
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Add the following directory to your PATH:")
	fmt.Println("    " + filepath.Join(homeDir, ".cvm", "bin"))
}
