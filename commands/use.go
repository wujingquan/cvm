package commands

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"wujingquan/cvm/common"
	"wujingquan/cvm/theme"
)

func Use(args []string) {
	if len(args) < 1 {
		theme.Error("You must specify a version to use.")
		return
	}

	// get users home dir
	homeDir, err := os.UserHomeDir()

	if err != nil {
		log.Fatalln(err)
	}

	// check if .cvm folder exists
	if _, err := os.Stat(filepath.Join(homeDir, ".cvm")); os.IsNotExist(err) {
		theme.Error("No Composer versions installed")
		return
	}

	// check if .cvm/versions folder exists
	if _, err := os.Stat(filepath.Join(homeDir, ".cvm", "versions")); os.IsNotExist(err) {
		theme.Error("No Composer versions installed")
		return
	}

	// check if .cvm/bin folder exists
	binPath := filepath.Join(homeDir, ".cvm", "bin")
	if _, err := os.Stat(binPath); os.IsNotExist(err) {
		os.Mkdir(binPath, 0755)
	}

	// get all folders in .cvm/versions
	versions, err := os.ReadDir(filepath.Join(homeDir, ".cvm", "versions"))
	if err != nil {
		log.Fatalln(err)
	}

	// transform to easily sortable slice
	var availableVersions []versionMeta
	for i, version := range versions {
		availableVersions = append(availableVersions, versionMeta{
			number: common.GetVersion(version.Name()),
			folder: versions[i],
		})
	}

	// check if version exists
	var selectedVersion *versionMeta
	for _, version := range availableVersions {
		if version.number.Major+"."+version.number.Minor+"."+version.number.Patch == args[0] {
			selectedVersion = &versionMeta{
				number: version.number,
				folder: version.folder,
			}
		}
	}

	// if patch version is not specified, use the newest matching major.minor
	if selectedVersion == nil {
		// Sort by newest patch first
		availableVersions = sortVersions(availableVersions)

		for _, version := range availableVersions {
			if version.number.Major+"."+version.number.Minor == args[0] {
				selectedVersion = &versionMeta{
					number: version.number,
					folder: version.folder,
				}
				break
			}
		}

		if selectedVersion == nil {
			theme.Error("The specified version is not installed.")
			return
		} else {
			theme.Warning(fmt.Sprintf("No patch version specified, assumed newest patch version %s.", selectedVersion.number.String()))
		}
	}

	// remove old bat script
	batPath := filepath.Join(binPath, "composer.bat")
	if _, err := os.Stat(batPath); err == nil {
		os.Remove(batPath)
	}

	// remove the old sh script
	shPath := filepath.Join(binPath, "composer")
	if _, err := os.Stat(shPath); err == nil {
		os.Remove(shPath)
	}

	versionFolderPath := filepath.Join(homeDir, ".cvm", "versions", selectedVersion.folder.Name())
	versionPath := filepath.Join(versionFolderPath, "composer.phar")

	// create bat script
	batCommand := "@echo off \n"
	batCommand = batCommand + ":: in case DelayedExpansion is on and a path contains ! \n"
	batCommand = batCommand + "setlocal DISABLEDELAYEDEXPANSION\n"
	batCommand = batCommand + "set filepath=\"" + versionPath + "\"\n"
	batCommand = batCommand + "set arguments=%*\n"
	batCommand = batCommand + "php %filepath% %arguments%\n"

	err = os.WriteFile(batPath, []byte(batCommand), 0755)

	if err != nil {
		log.Fatalln(err)
	}

	// create sh script
	shCommand := "#!/bin/bash\n"
	shCommand = shCommand + "filepath=\"" + versionPath + "\"\n"
	shCommand = shCommand + "\"$filepath\" \"$@\""

	err = os.WriteFile(shPath, []byte(shCommand), 0755)

	if err != nil {
		log.Fatalln(err)
	}

	theme.Success("Using Composer " + selectedVersion.number.String())
}

func sortVersions(in []versionMeta) []versionMeta {
	sort.Slice(in, func(i, j int) bool {
		if in[i].number.Major != in[j].number.Major {
			return in[i].number.Major > in[j].number.Major
		}
		if in[i].number.Minor != in[j].number.Minor {
			return in[i].number.Minor > in[j].number.Minor
		}
		return in[i].number.Patch > in[j].number.Patch
	})

	return in
}

type versionMeta struct {
	number common.Version
	folder os.DirEntry
}
