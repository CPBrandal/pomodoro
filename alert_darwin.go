//go:build darwin

package main

import (
	"fmt"
	"os"
	"os/exec"
)

func alert(message string) {
	cmd := exec.Command("osascript", "-e", fmt.Sprintf(`display dialog "%s" with title "P O M O D O R O"`, message))
	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			fmt.Println("Pomodoro cancelled.")
			os.Exit(0)
		}
		fmt.Println("Failed to send notification:", err)
	}
}