package main

import (
	"bufio"
	"embed"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var lastPreset LastPresetChoice
var asciiPaintings embed.FS

func main() {	
	setupExitHandler()
	
	if len(os.Args) > 1 {
		arg := strings.ToLower(os.Args[1])
		switch arg {
		case "--uninstall", "-u", "uninstall":
			uninstall()
			return
		case "--help", "-h", "help":
			showCLIHelp()
			return
		default:
			fmt.Printf("Unknown option: %s\n", os.Args[1])
			fmt.Println("Use 'pomodoro --help' for usage information")
			os.Exit(1)
			return
		}
	}

	presets := loadPresets()
	if presets.LastUsedPreset != nil {
		lastPreset = LastPresetChoice{
			WorkTime:  time.Duration(presets.LastUsedPreset.WorkMinutes) * time.Minute,
			BreakTime: time.Duration(presets.LastUsedPreset.BreakMinutes) * time.Minute,
			LongerBreakTime: time.Duration(presets.LastUsedPreset.LongerBreakMinutes) * time.Minute,
		}
	} else {
		lastPreset = LastPresetChoice{
			WorkTime:  time.Duration(DEFAULT_WORK_DURATION) * time.Minute,
			BreakTime: time.Duration(DEFAULT_BREAK_DURATION) * time.Minute,
			LongerBreakTime: time.Duration(DEFAULT_LONGER_BREAK_DURATION) * time.Minute,
		}
	}
	showMainMenu()
}

func setupExitHandler() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		os.Exit(0)
	}()
}

func showMainMenu() {	
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\n================================================")
	fmt.Printf("%sP O M O D O R O%s\n", "\033[31m", "\033[0m")
	fmt.Println("================================================")
	fmt.Print("1 - Use default values (25 min work, 5 min break)\n")
	fmt.Print("2 - Select custom values\n")
	fmt.Print("3 - Delete presets\n\n")
	fmt.Print("a - Artwork Gallery\n")
	fmt.Print("⏎ - Use last selected preset\n")
	fmt.Print("q - Quit the program\n")
	printUserInputPrompt()
	for {
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)
		
		switch choice {
		case "q":
			os.Exit(0)
			return
		case "1":
			workDuration := time.Duration(DEFAULT_WORK_DURATION) * time.Minute
			breakDuration := time.Duration(DEFAULT_BREAK_DURATION) * time.Minute
			longerBreakDuration := time.Duration(DEFAULT_LONGER_BREAK_DURATION) * time.Minute
			workBreakLoop(workDuration, breakDuration, longerBreakDuration)
			return
		case "2":
			selectCustomValues(reader)
			return
		case "3":
			deletePresets(reader)
			return
		case "":
            	if lastPreset.WorkTime > 0 && lastPreset.BreakTime > 0 {
				fmt.Printf("\nUsing last selected preset, %.0f min work, %.0f min break, %.0f min longer break\n", lastPreset.WorkTime.Minutes(), lastPreset.BreakTime.Minutes(), lastPreset.LongerBreakTime.Minutes())
                workBreakLoop(
                    lastPreset.WorkTime,
                    lastPreset.BreakTime,
                    lastPreset.LongerBreakTime)
                return
            }
            fmt.Println("\nPlease select a custom timer first.")
			printUserInputPrompt()
		case "a":
			showArtworkGallery(reader)
			printUserInputPrompt()
		case "h":
			showHelp()
			printUserInputPrompt()
		default:
			fmt.Println("\nNot sure what that is. Try 'h' for help.")
			printUserInputPrompt()
		}
	}
}

func showHelp() {
	fmt.Println("\nCommands:")
	fmt.Println("1 ─ Default timer (25/5 min)")
	fmt.Println("2 ─ Presets & custom timers")
	fmt.Println("3 ─ Delete presets")
	fmt.Println("a ─ Artwork Gallery")
	fmt.Println("⏎ (Enter) ─ Last custom timer")
	fmt.Println("\nh ─ Help")
	fmt.Println("q ─ Quit")
}

func selectCustomValues(reader *bufio.Reader) {
	presets := loadPresets()
	
	if len(presets.Presets) > 0 {
		fmt.Println("\nSaved presets:")
		for i, preset := range presets.Presets {
			fmt.Printf("%d - %s (%d min work, %d min break)\n", i+1, preset.Name, preset.WorkMinutes, preset.BreakMinutes)
		}
		fmt.Printf("%d - Create new custom timer\n", len(presets.Presets)+1)
		fmt.Printf("%d - Return to main menu\n", 0)
		printUserInputPrompt()
		
		choiceStr, _ := reader.ReadString('\n')
		choiceStr = strings.TrimSpace(choiceStr)
		choice, err := strconv.Atoi(choiceStr)
		if err == nil && choice == 0 {
			showMainMenu()
			return
		} else if err == nil && choice >= 1 && choice <= len(presets.Presets) {
			selected := presets.Presets[choice-1]
			workDuration := time.Duration(selected.WorkMinutes) * time.Minute
			breakDuration := time.Duration(selected.BreakMinutes) * time.Minute
			longerBreakDuration := time.Duration(selected.LongerBreakMinutes) * time.Minute
			
			saveLastUsedPreset(selected.WorkMinutes, selected.BreakMinutes, selected.LongerBreakMinutes)
			lastPreset = LastPresetChoice{WorkTime: workDuration, BreakTime: breakDuration, LongerBreakTime: longerBreakDuration}
			
			workBreakLoop(workDuration, breakDuration, longerBreakDuration)
			return
		} else if err == nil && choice == len(presets.Presets)+1 {
			createCustomTimer(reader)
			return
		} else {
			fmt.Println("Invalid selection.")
			return
		}
	} else {
		createCustomTimer(reader)
	}
}

func createCustomTimer(reader *bufio.Reader) {
	workMinutes := userInputHandler(reader, "How long are your working intervals (in minutes) [25]: ", DEFAULT_WORK_DURATION)
	if workMinutes < 0 {
		return
	}
	
	breakMinutes := userInputHandler(reader, "How long is your break for (in minutes) [5]: ", DEFAULT_BREAK_DURATION)
	if breakMinutes < 0 {
		return
	}
	
	longerBreakMinutes := userInputHandler(reader, "How long is your longer break for (in minutes) [20]: ", DEFAULT_LONGER_BREAK_DURATION)
	if longerBreakMinutes < 0 {
		return
	}
	
	fmt.Print("Save as preset? (y/n) [n]: ")
	saveStr, _ := reader.ReadString('\n')
	saveStr = strings.TrimSpace(strings.ToLower(saveStr))
	if saveStr == "y" || saveStr == "yes" {
		fmt.Print("Preset name: ")
		nameStr, _ := reader.ReadString('\n')
		nameStr = strings.TrimSpace(nameStr)
		if nameStr != "" {
			savePreset(nameStr, workMinutes, breakMinutes, longerBreakMinutes)
			fmt.Printf("Preset '%s' saved!\n", nameStr)
		}
	}
	
	workDuration := time.Duration(workMinutes) * time.Minute
	breakDuration := time.Duration(breakMinutes) * time.Minute
	longerBreakDuration := time.Duration(longerBreakMinutes) * time.Minute

	saveLastUsedPreset(workMinutes, breakMinutes, longerBreakMinutes)
	lastPreset = LastPresetChoice{WorkTime: workDuration, BreakTime: breakDuration, LongerBreakTime: longerBreakDuration}

	workBreakLoop(workDuration, breakDuration, longerBreakDuration)
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

func workBreakLoop(workDuration time.Duration, breakDuration time.Duration, longerBreakDuration time.Duration) {	
	saveLastUsedPreset(int(workDuration.Minutes()), int(breakDuration.Minutes()), int(longerBreakDuration.Minutes()))
	
	for i := range 4 {
		fmt.Printf("Work session %d started...\n", i+1)
		time.Sleep(workDuration)
		alert(fmt.Sprintf("Take a break! You worked for %.0f minutes.\nA %.0f minute break starts now.", workDuration.Minutes(), breakDuration.Minutes()))
		addSessionTime(workDuration, 0)
		if(i==3) {break;}
		fmt.Printf("Break time (%.0f minutes)...\n", breakDuration.Minutes())
		time.Sleep(breakDuration)
		alert("Break over! Time to get back to work.")
	}

	unlockArtworkLines(4)
	progress := loadArtworkProgress()
	if progress.CurrentArtworkIndex >= len(artworkList) {
		fmt.Println("\nCongratulations! You've completed all artworks!")
	} else {
		currentArtwork := artworkList[progress.CurrentArtworkIndex]
		unlockedLines := progress.UnlockedLines[currentArtwork.Filename]
		totalArtLines := currentArtwork.TotalLines - 1
		
		fmt.Println("\nArtwork Progress:")
		fmt.Printf("   %s: %d/%d lines unlocked\n", currentArtwork.Name, unlockedLines, totalArtLines)
		
		if unlockedLines > 0 {
			fmt.Println("\n   Preview:")
			displayPartialArtwork(currentArtwork.Filename, unlockedLines)
		}
		
		if unlockedLines >= totalArtLines {
			fmt.Printf("\n   Congratulations! You've completed %s!\n", currentArtwork.Name)
		}
	}

	alert(fmt.Sprintf("Great job! Time for a longer %.0f minute break.", longerBreakDuration.Minutes()))
	time.Sleep(longerBreakDuration)
	alert("You have completed your pomodoro session. Press ok to restart, or cancel to exit.")
	showMainMenu()
}

func printUserInputPrompt() {
	fmt.Print("\n❯  ")
}
