package commands

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"wujingquan/cvm/common"
	"wujingquan/cvm/theme"
)

type Version struct {
	Major string
	Minor string
	Patch string
	Url   string
}

func Install(args []string) {
	if len(args) < 2 {
		theme.Error("You must specify a version to install.")
		return
	}

	desiredVersionNumbers := common.GetVersion(args[1])

	if desiredVersionNumbers == (common.Version{}) {
		theme.Error("Invalid version specified")
		return
	}

	// Get the desired version from the user input
	desiredMajorVersion := desiredVersionNumbers.Major
	desiredMinorVersion := desiredVersionNumbers.Minor
	desiredPatchVersion := desiredVersionNumbers.Patch

	// perform get request to https://getcomposer.org/download/
	resp, err := http.Get("https://getcomposer.org/download/")
	if err != nil {
		log.Fatalln(err)
	}
	// We Read the response body on the line below.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	// Convert the body to type string
	sb := string(body)
	re := regexp.MustCompile(`<a href="([a-zA-Z0-9./-]+composer.phar)"([a-zA-Z0-9\s./-]+)title="([a-zA-Z0-9\s./-]+)"([a-zA-Z0-9\s./-]+)aria-label="([a-zA-Z0-9\s./-]+)">([a-zA-Z0-9./-]+)<\/a>`)
	matches := re.FindAllStringSubmatch(sb, -1)

	versions := make([]Version, 0)

	for _, match := range matches {
		url := match[1]
		name := match[6]

		// regex match name
		versionNumbers := common.GetVersion(name)

		major := versionNumbers.Major
		minor := versionNumbers.Minor
		patch := versionNumbers.Patch

		// push to versions
		versions = append(versions, Version{
			Major: major,
			Minor: minor,
			Patch: patch,
			Url:   url,
		})
	}

	// find desired version
	var desiredVersion Version

	if desiredMajorVersion != "" && desiredMinorVersion != "" && desiredPatchVersion != "" {
		desiredVersion = FindExactVersion(versions, desiredMajorVersion, desiredMinorVersion, desiredPatchVersion)
	}

	if desiredMajorVersion != "" && desiredMinorVersion != "" && desiredPatchVersion == "" {
		desiredVersion = FindLatestPatch(versions, desiredMajorVersion, desiredMinorVersion)
	}

	if desiredMajorVersion != "" && desiredMinorVersion == "" && desiredPatchVersion == "" {
		desiredVersion = FindLatestMinor(versions, desiredMajorVersion)
	}

	if desiredVersion == (Version{}) {
		theme.Error("Could not find the desired version: " + args[1])
		return
	}

	fmt.Println("Installing Composer " + desiredVersion.Major + "." + desiredVersion.Minor + "." + desiredVersion.Patch)

	homeDir, err := os.UserHomeDir()

	if err != nil {
		log.Fatalln(err)
	}

	// check if .cvm folder exists
	if _, err := os.Stat(homeDir + "/.cvm"); os.IsNotExist(err) {
		theme.Info("Creating .cvm folder in home directory")
		os.Mkdir(homeDir+"/.cvm", 0755)
	}

	// check if .cvm/versions folder exists
	if _, err := os.Stat(homeDir + "/.cvm/versions"); os.IsNotExist(err) {
		theme.Info("Creating .cvm/versions folder in home directory")
		os.Mkdir(homeDir+"/.cvm/versions", 0755)
	}

	// check if .cvm/versions/[specified version path] folder exists
	if _, err := os.Stat(homeDir + "/.cvm/versions/" + desiredVersion.Major + "." + desiredVersion.Minor + "." + desiredVersion.Patch); os.IsNotExist(err) {
		theme.Info("Creating .cvm/versions/" + desiredVersion.Major + "." + desiredVersion.Minor + "." + desiredVersion.Patch + " folder in home directory")
		os.Mkdir(homeDir+"/.cvm/versions/"+desiredVersion.Major+"."+desiredVersion.Minor+"."+desiredVersion.Patch, 0755)
	}

	theme.Info("Downloading")

	// Get the data
	downloadResponse, err := http.Get("https://getcomposer.org" + desiredVersion.Url)
	if err != nil {
		log.Fatalln(err)
	}

	defer downloadResponse.Body.Close()

	// filename from url
	filename := strings.Split(desiredVersion.Url, "/")[len(strings.Split(desiredVersion.Url, "/"))-1]

	// check if composer already exists
	if _, err := os.Stat(homeDir + "/.cvm/versions/" + desiredVersion.Major + "." + desiredVersion.Minor + "." + desiredVersion.Patch + "/" + filename); err == nil {
		theme.Error("Composer " + desiredVersion.Major + "." + desiredVersion.Minor + "." + desiredVersion.Patch + " already exists")
		return
	}

	// Create the file
	out, err := os.Create(homeDir + "/.cvm/versions/" + desiredVersion.Major + "." + desiredVersion.Minor + "." + desiredVersion.Patch + "/" + filename)
	if err != nil {
		log.Fatalln(err)
	}

	// Write the body to file
	_, err = io.Copy(out, downloadResponse.Body)

	if err != nil {
		out.Close()
		log.Fatalln(err)
	}

	// Close the file
	out.Close()

	theme.Success("Finished installing Composer " + desiredVersion.Major + "." + desiredVersion.Minor + "." + desiredVersion.Patch)
}

func FindExactVersion(versions []Version, major string, minor string, patch string) Version {
	for _, version := range versions {
		if version.Major == major && version.Minor == minor && version.Patch == patch {
			return version
		}
	}

	return Version{}
}

func FindLatestPatch(versions []Version, major string, minor string) Version {
	latestPatch := Version{}

	for _, version := range versions {
		if version.Major == major && version.Minor == minor {
			if latestPatch.Patch == "" || version.Patch > latestPatch.Patch {
				latestPatch = version
			}
		}
	}

	return latestPatch
}

func FindLatestMinor(versions []Version, major string) Version {
	latestMinor := Version{}

	for _, version := range versions {
		if version.Major == major {
			if latestMinor.Minor == "" || version.Minor > latestMinor.Minor {
				if latestMinor.Patch == "" || version.Patch > latestMinor.Patch {
					latestMinor = version
				}
			}
		}
	}

	return latestMinor
}
