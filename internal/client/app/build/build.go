package build

import "fmt"

var (
	Version string // global buildVersion
	Date    string // global buildDate
	Commit  string // global buildCommit
)

const notGiven = "N/A"

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

func PrintBuildInfo() {
	fmt.Printf("Build version: %s\n", Version) // print build info
	fmt.Printf("Build date: %s\n", Date)       // print build info
	fmt.Printf("Build commit: %s\n", Commit)   // print build info
}
