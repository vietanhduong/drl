package config

import (
	"fmt"
	"os"
)

var (
	version   = "unreleased"
	gitCommit = "none"
	buildDate = "unknown"
)

func PrintVersion() {
	fmt.Fprintf(os.Stdout, "Version: %s\n", version)
	fmt.Fprintf(os.Stdout, "Git commit: %s\n", gitCommit)
	fmt.Fprintf(os.Stdout, "Build date: %s\n", buildDate)
}
