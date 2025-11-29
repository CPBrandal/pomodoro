//go:build linux

package main

import (
	"fmt"
	"os/exec"
)

func alert(message string) {
	cmd := exec.Command("notify-send", "Pomodoro", message)
	if err := cmd.Run(); err != nil {
		fmt.Println("Failed to send notification:", err)
	}
}