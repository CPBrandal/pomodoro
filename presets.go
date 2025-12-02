package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

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

func saveLastUsedPreset(workMinutes, breakMinutes, longerBreakMinutes int) {
	presets := loadPresets()
	
	presets.LastUsedPreset = &Preset{
		Name:         "Last Used",
		WorkMinutes:  workMinutes,
		BreakMinutes: breakMinutes,
		LongerBreakMinutes: longerBreakMinutes,
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

func savePreset(name string, workMinutes, breakMinutes, longerBreakMinutes int) {
	presets := loadPresets()
	
	found := false
	for i := range presets.Presets {
		if presets.Presets[i].Name == name {
			presets.Presets[i].WorkMinutes = workMinutes
			presets.Presets[i].BreakMinutes = breakMinutes
			presets.Presets[i].LongerBreakMinutes = longerBreakMinutes
			found = true
			break
		}
	}
	
	if !found {
		presets.Presets = append(presets.Presets, Preset{
			Name:         name,
			WorkMinutes:  workMinutes,
			BreakMinutes: breakMinutes,
			LongerBreakMinutes: longerBreakMinutes,
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

func deletePresets(reader *bufio.Reader) {
	for {
		presets := loadPresets()
		
		if len(presets.Presets) == 0 {
			fmt.Println("\nNo presets to delete.")
			showMainMenu()
			return
		}
		
		fmt.Println("\nDelete a preset:")
		for i, preset := range presets.Presets {
			fmt.Printf("%d - %s (%d min work, %d min break)\n", i+1, preset.Name, preset.WorkMinutes, preset.BreakMinutes)
		}
		fmt.Printf("%d - Return to main menu\n", 0)
		printUserInputPrompt()
			choiceStr, _ := reader.ReadString('\n')
			choiceStr = strings.TrimSpace(choiceStr)
			choice, err := strconv.Atoi(choiceStr)

			if err != nil || choice < 0 || choice > len(presets.Presets) {
				fmt.Println("\nInvalid input. Enter a number between 1 and", len(presets.Presets), "to delete a preset, or 0 to return to main menu.")
				continue
			}

			if choice == 0 {
				showMainMenu()
				return
			}

			presetToDelete := presets.Presets[choice-1]
			presets.Presets = append(presets.Presets[:choice-1], presets.Presets[choice:]...)

			configPath, err := getConfigPath()
			if err != nil {
				fmt.Println("Error: Could not save changes.")
				continue
			}

			data, err := json.MarshalIndent(presets, "", "  ")
			if err != nil {
				fmt.Println("Error: Could not save changes.")
				continue
			}

			os.WriteFile(configPath, data, 0644)
			fmt.Printf("Preset '%s' deleted successfully.\n", presetToDelete.Name)
		}
}