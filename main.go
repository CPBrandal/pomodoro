package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func alert(message string) {
	cmd := exec.Command("osascript", "-e", fmt.Sprintf(`display dialog "%s" with title "Pomodoro"`, message))
	err := cmd.Run()
	if err != nil {
		fmt.Println("Failed to send notification:", err)
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("How long are your working intervals (in minutes) [25]: ")
	durationStr, _ := reader.ReadString('\n')
	durationStr = strings.TrimSpace(durationStr)

	var workMinutes int
	if durationStr == "" {
		workMinutes = 25
	} else {
		var err error
		workMinutes, err = strconv.Atoi(durationStr)
		if err != nil || workMinutes <= 0 {
			fmt.Println("Please enter a valid positive number for duration.")
			return
		}
	}

	fmt.Print("How long is your break for (in minutes) [5]: ")
	breakStr, _ := reader.ReadString('\n')
	breakStr = strings.TrimSpace(breakStr)

	var breakMinutes int
	if breakStr == "" {
		breakMinutes = 5
	} else {
		var err error
		breakMinutes, err = strconv.Atoi(breakStr)
		if err != nil || breakMinutes <= 0 {
			fmt.Println("Please enter a valid positive number for duration.")
			return
		}
	}

	workDuration := time.Duration(workMinutes) * time.Minute
	breakDuration := time.Duration(breakMinutes) * time.Minute
	
	workLoop(workDuration, breakDuration)
}

func workLoop(workDuration time.Duration, breakDuration time.Duration) {
	for i := range 4 {
		fmt.Printf("Work session %d started...\n", i+1)
		time.Sleep(workDuration)
		alert(fmt.Sprintf("Take a break! You worked for %.0f minutes. \n A %.0f minute break starts now.", workDuration.Minutes(), breakDuration.Minutes()))

		fmt.Printf("Break time (%.0f minutes)...\n", breakDuration.Minutes())
		time.Sleep(breakDuration)
		alert("Break over! Time to get back to work.")
	}

	alert("Great job! Time for a longer 20 minute break.")
	time.Sleep(20 * time.Minute)
}