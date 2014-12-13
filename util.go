package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os/exec"

	"github.com/mgutz/ansi"
)

func printError(message string, args ...interface{}) {
	log.Println(colorizeMessage("red", "error:", message, args...))
}

func printFatal(message string, args ...interface{}) {
	log.Fatal(colorizeMessage("red", "error:", message, args...))
}

func printWarning(message string, args ...interface{}) {
	log.Println(colorizeMessage("yellow", "warning:", message, args...))
}

func colorizeMessage(color, prefix, message string, args ...interface{}) string {
	prefResult := ""
	if prefix != "" {
		prefResult = ansi.Color(prefix, color+"+b") + " " + ansi.ColorCode("reset")
	}
	return prefResult + ansi.Color(fmt.Sprintf(message, args...), color) + ansi.ColorCode("reset")
}

func listRec(w io.Writer, a ...interface{}) {
	for i, x := range a {
		fmt.Fprint(w, x)
		if i+1 < len(a) {
			w.Write([]byte{'\t'})
		} else {
			w.Write([]byte{'\n'})
		}
	}
}

func runCommand(cmd string) (string, string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	c := exec.Command("/bin/bash", "-lc", cmd)
	c.Stdout = &stdout
	c.Stderr = &stderr

	err := c.Run()
	return string(stdout.Bytes()), string(stderr.Bytes()), err
}
