package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const DEFAULT_WORK_DURATION = 25
const DEFAULT_BREAK_DURATION = 5

type Preset struct {
	Name         string `json:"name"`
	WorkMinutes  int    `json:"workMinutes"`
	BreakMinutes int    `json:"breakMinutes"`
}

type Presets struct {
	Presets []Preset `json:"presets"`
}

func main() {
	startPomodoro(true)
}

func startPomodoro(promptUser bool) {
	if !promptUser {
		workDuration := time.Duration(DEFAULT_WORK_DURATION) * time.Minute
		breakDuration := time.Duration(DEFAULT_BREAK_DURATION) * time.Minute
		workBreakLoop(workDuration, breakDuration)
		return
	}

	reader := bufio.NewReader(os.Stdin)

	// Show preset menu
	presets := loadPresets()
	if len(presets.Presets) > 0 {
		fmt.Println("\nSaved presets:")
		for i, preset := range presets.Presets {
			fmt.Printf("  %d. %s (%d min work, %d min break)\n", i+1, preset.Name, preset.WorkMinutes, preset.BreakMinutes)
		}
		fmt.Printf("  %d. Create new preset\n", len(presets.Presets)+1)
		fmt.Printf("  %d. Use default (25/5)\n", len(presets.Presets)+2)
		fmt.Print("\nSelect an option: ")

		choiceStr, _ := reader.ReadString('\n')
		choiceStr = strings.TrimSpace(choiceStr)
		choice, err := strconv.Atoi(choiceStr)

		if err == nil && choice >= 1 && choice <= len(presets.Presets) {
			// Use selected preset
			selected := presets.Presets[choice-1]
			workDuration := time.Duration(selected.WorkMinutes) * time.Minute
			breakDuration := time.Duration(selected.BreakMinutes) * time.Minute
			workBreakLoop(workDuration, breakDuration)
			return
		} else if err == nil && choice == len(presets.Presets)+2 {
			// Use default
			workDuration := time.Duration(DEFAULT_WORK_DURATION) * time.Minute
			breakDuration := time.Duration(DEFAULT_BREAK_DURATION) * time.Minute
			workBreakLoop(workDuration, breakDuration)
			return
		}
		// Otherwise continue to create new preset
	}

	// Get user input for new preset
	workMinutes := userInputHandler(reader, "How long are your working intervals (in minutes) [25]: ", DEFAULT_WORK_DURATION)
	if workMinutes < 0 {
		return
	}

	breakMinutes := userInputHandler(reader, "How long is your break for (in minutes) [5]: ", DEFAULT_BREAK_DURATION)
	if breakMinutes < 0 {
		return
	}

	// Ask if user wants to save as preset
	fmt.Print("Save as preset? (y/n) [n]: ")
	saveStr, _ := reader.ReadString('\n')
	saveStr = strings.TrimSpace(strings.ToLower(saveStr))
	if saveStr == "y" || saveStr == "yes" {
		fmt.Print("Preset name: ")
		nameStr, _ := reader.ReadString('\n')
		nameStr = strings.TrimSpace(nameStr)
		if nameStr != "" {
			savePreset(nameStr, workMinutes, breakMinutes)
			fmt.Printf("Preset '%s' saved!\n", nameStr)
		}
	}

	workDuration := time.Duration(workMinutes) * time.Minute
	breakDuration := time.Duration(breakMinutes) * time.Minute

	workBreakLoop(workDuration, breakDuration)
}

func userInputHandler(reader *bufio.Reader, msg string, defaultValue int) int {
	fmt.Print(msg)
	durationStr, _ := reader.ReadString('\n')
	durationStr = strings.TrimSpace(durationStr)

	if durationStr == "" {
		return defaultValue
	}

	minutes, err := strconv.Atoi(durationStr)
	if err != nil || minutes <= 0 {
		fmt.Println("Please enter a valid positive number for duration.")
		return -1
	}
	return minutes
}

func getConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	configDir := filepath.Join(homeDir, ".pomodoro")
	os.MkdirAll(configDir, 0755)
	return filepath.Join(configDir, "presets.json"), nil
}

func loadPresets() Presets {
	configPath, err := getConfigPath()
	if err != nil {
		return Presets{Presets: []Preset{}}
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return Presets{Presets: []Preset{}}
	}

	var presets Presets
	if err := json.Unmarshal(data, &presets); err != nil {
		return Presets{Presets: []Preset{}}
	}

	return presets
}

func savePreset(name string, workMinutes, breakMinutes int) {
	presets := loadPresets()
	
	// Check if preset with same name exists and update it
	found := false
	for i := range presets.Presets {
		if presets.Presets[i].Name == name {
			presets.Presets[i].WorkMinutes = workMinutes
			presets.Presets[i].BreakMinutes = breakMinutes
			found = true
			break
		}
	}
	
	// Add new preset if not found
	if !found {
		presets.Presets = append(presets.Presets, Preset{
			Name:         name,
			WorkMinutes:  workMinutes,
			BreakMinutes: breakMinutes,
		})
	}

	configPath, err := getConfigPath()
	if err != nil {
		return
	}

	data, err := json.MarshalIndent(presets, "", "  ")
	if err != nil {
		return
	}

	os.WriteFile(configPath, data, 0644)
}

func workBreakLoop(workDuration time.Duration, breakDuration time.Duration) {
	for i := range 4 {
		fmt.Printf("Work session %d started...\n", i+1)
		//time.Sleep(workDuration)
		alert(fmt.Sprintf("Take a break! You worked for %.0f minutes.\nA %.0f minute break starts now.", workDuration.Minutes(), breakDuration.Minutes()))

		if(i!=3) {
			fmt.Printf("Break time (%.0f minutes)...\n", breakDuration.Minutes())
			//time.Sleep(breakDuration)
			alert("Break over! Time to get back to work.")
		}
	}

	alert("Great job! Time for a longer 20 minute break.")
	//time.Sleep(20 * time.Minute)
	alert("You have completed your pomodoro session. Press ok to restart, or cancel to exit.")
	startPomodoro(false)
}