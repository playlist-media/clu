package main

import (
	"os/exec"
	"strings"
)

func gitConfigBool(name string) bool {
	b, err := exec.Command("git", "config", name).Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(b)) == "true"
}
