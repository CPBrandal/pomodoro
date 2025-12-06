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

func getUsagePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	configDir := filepath.Join(homeDir, ".pomodoro")
	os.MkdirAll(configDir, 0755)
	return filepath.Join(configDir, "usage.json"), nil
}

func loadUsageStats() UsageStats {
	usagePath, err := getUsagePath()
	if err != nil {
		return UsageStats{TotalHours: 0, TotalPomodoros: 0}
	}

	data, err := os.ReadFile(usagePath)
	if err != nil {
		return UsageStats{TotalHours: 0, TotalPomodoros: 0}
	}

	var stats UsageStats
	if err := json.Unmarshal(data, &stats); err != nil {
		return UsageStats{TotalHours: 0, TotalPomodoros: 0}
	}

	return stats
}

func saveUsageStats(stats UsageStats) error {
	usagePath, err := getUsagePath()
	if err != nil {
		return err
	}

	stats.LastUpdated = time.Now()
	data, err := json.MarshalIndent(stats, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(usagePath, data, 0644)
}

func addSessionTime(duration time.Duration, pomodoros int) error {
	stats := loadUsageStats()
	stats.TotalHours += duration.Hours()
	stats.TotalPomodoros += pomodoros
	return saveUsageStats(stats)
}

func getTotalUsageHours() float64 {
	stats := loadUsageStats()
	return stats.TotalHours
}

func getTotalPomodoros() int {
	stats := loadUsageStats()
	return stats.TotalPomodoros
}

func getArtworkProgressPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	configDir := filepath.Join(homeDir, ".pomodoro")
	os.MkdirAll(configDir, 0755)
	return filepath.Join(configDir, "artwork.json"), nil
}

func loadArtworkProgress() ArtworkProgress {
	progressPath, err := getArtworkProgressPath()
	if err != nil {
		return ArtworkProgress{
			CurrentArtworkIndex: 0,
			UnlockedLines:       make(map[string]int),
		}
	}

	data, err := os.ReadFile(progressPath)
	if err != nil {
		return ArtworkProgress{
			CurrentArtworkIndex: 0,
			UnlockedLines:       make(map[string]int),
		}
	}

	var progress ArtworkProgress
	if err := json.Unmarshal(data, &progress); err != nil {
		return ArtworkProgress{
			CurrentArtworkIndex: 0,
			UnlockedLines:       make(map[string]int),
		}
	}

	if progress.UnlockedLines == nil {
		progress.UnlockedLines = make(map[string]int)
	}

	return progress
}

func saveArtworkProgress(progress ArtworkProgress) error {
	progressPath, err := getArtworkProgressPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(progress, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(progressPath, data, 0644)
}

func unlockArtworkLines(lines int) {
	progress := loadArtworkProgress()
	
	if progress.CurrentArtworkIndex >= len(artworkList) {
		return
	}

	currentArtwork := artworkList[progress.CurrentArtworkIndex]
	currentUnlocked := progress.UnlockedLines[currentArtwork.Filename]
	
	newUnlocked := currentUnlocked + lines
	totalArtLines := currentArtwork.TotalLines - 1
	
	if newUnlocked >= totalArtLines {
		progress.UnlockedLines[currentArtwork.Filename] = totalArtLines
		progress.CurrentArtworkIndex++
	} else {
		progress.UnlockedLines[currentArtwork.Filename] = newUnlocked
	}
	
	saveArtworkProgress(progress)
}

func displayPartialArtwork(filename string, linesToShow int) error {
	file, err := asciiPaintings.Open(filename)
	if err != nil {
		file, err = os.Open(filename)
		if err != nil {
			return err
		}
	}
	defer file.Close()
	
	scanner := bufio.NewScanner(file)
	lineNum := 0
	artLinesShown := 0
	
	for scanner.Scan() {
		line := scanner.Text()
		lineNum++
		
		// Skip DIMENSIONS line
		if strings.HasPrefix(line, "DIMENSIONS:") {
			continue
		}
		
		// Show art lines up to linesToShow
		if artLinesShown < linesToShow {
			fmt.Println(line)
			artLinesShown++
		} else {
			break
		}
	}
	
	return scanner.Err()
}

func printArtworkGalleryHeader() {
	fmt.Println("\n================================================")
	fmt.Printf("%sA R T W O R K   G A L L E R Y%s\n", "\033[31m", "\033[0m")
	fmt.Println("================================================")
}

func showArtworkGallery(reader *bufio.Reader) {
	for {
		progress := loadArtworkProgress()
		totalHours := getTotalUsageHours()
		totalPomodoros := getTotalPomodoros()
		hours := int(totalHours)
		minutes := int((totalHours - float64(hours)) * 60)
		
		printArtworkGalleryHeader()
		fmt.Printf("Total worktime: %d hours, %d minutes (%.2f hours)\n", hours, minutes, totalHours)
		fmt.Printf("Total pomodoro sessions completed: %d\n\n", totalPomodoros)
		
		for i, artwork := range artworkList {
			unlockedLines := progress.UnlockedLines[artwork.Filename]
			totalArtLines := artwork.TotalLines - 1
			status := ""
			
			if i < progress.CurrentArtworkIndex {
				status = "Completed"
			} else if i == progress.CurrentArtworkIndex {
				if unlockedLines >= totalArtLines {
					status = "Completed"
				} else {
					status = fmt.Sprintf("In Progress (%d/%d lines)", unlockedLines, totalArtLines)
				}
			} else {
				status = "Locked"
			}
			
			fmt.Printf("%d - %s - %s\n", i+1, artwork.Name, status)
		}
		
		fmt.Println("\n0 - Return to main menu")
		printUserInputPrompt()
		
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)
		
		if choice == "0" {
			showMainMenu()
			return
		}
		
		choiceNum, err := strconv.Atoi(choice)
		if err == nil && choiceNum >= 1 && choiceNum <= len(artworkList) {
			artworkIndex := choiceNum - 1
			artwork := artworkList[artworkIndex]
			unlockedLines := progress.UnlockedLines[artwork.Filename]
			
			if artworkIndex < progress.CurrentArtworkIndex {
				showASCIIArtWithAutoSize(artwork.Filename)
			} else if artworkIndex == progress.CurrentArtworkIndex && unlockedLines > 0 {
				fmt.Printf("\n%s - Progress (%d/%d lines):\n\n", artwork.Name, unlockedLines, artwork.TotalLines-1)
				displayPartialArtwork(artwork.Filename, unlockedLines)
			} else {
				fmt.Printf("\n%s is locked. Complete previous artworks to unlock it.\n", artwork.Name)
				printUserInputPrompt()
				reader.ReadString('\n')
			}
			printUserInputPrompt()
			reader.ReadString('\n')
			continue
		} else {
			fmt.Println("Invalid selection. Enter a number between 1 and", len(artworkList), "or 0 to return.")
			printUserInputPrompt()
		}
	}
}
