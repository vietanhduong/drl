package config

import (
	"fmt"
	"os"
)

var (
	version   = "unreleased"
	gitCommit = "none"
	buildDate = "unknown"
	arch      = "unknown"
)

func PrintVersion() {
	fmt.Fprintf(os.Stdout, "Version:    %s\n", version)
	fmt.Fprintf(os.Stdout, "Git Commit: %s\n", gitCommit)
	fmt.Fprintf(os.Stdout, "Build Date: %s\n", buildDate)
	fmt.Fprintf(os.Stdout, "Arch:       %s\n", arch)
}
