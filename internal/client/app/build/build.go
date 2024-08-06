package build

import (
	"fmt"
)

var (
	Version string // Version represents the global build version.
	Date    string // Date represents the global build date.
	Commit  string // Commit represents the global build commit hash.
)

const notGiven = "N/A"

// CheckBuild verifies if the global build variables are set.
// If any of them are not set, it assigns the constant value "N/A".
func CheckBuild() {
	if Version == "" {
		Version = notGiven
	}

	if Date == "" {
		Date = notGiven
	}

	if Commit == "" {
		Commit = notGiven
	}
}

// PrintBuildInfo prints the current build information to the standard output.
// It displays the build version, date, and commit hash.
func PrintBuildInfo() {
	fmt.Printf("Build version: %s\n", Version)
	fmt.Printf("Build date: %s\n", Date)
	fmt.Printf("Build commit: %s\n", Commit)
}
