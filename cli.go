package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func showCLIHelp() {
	fmt.Println("Pomodoro Timer")
	fmt.Println("\nUsage:")
	fmt.Println("  pomodoro              Start the timer")
	fmt.Println("  pomodoro --uninstall  Uninstall the program and remove config")
	fmt.Println("  pomodoro -u          Short form for uninstall")
	fmt.Println("  pomodoro --stats     Show total usage statistics")
	fmt.Println("  pomodoro -s          Short form for stats")
	fmt.Println("  pomodoro --help      Show this help message")
	fmt.Println("  pomodoro -h          Short form for help")
}

func uninstall() {
	fmt.Println("Uninstalling Pomodoro Timer...")
	
	installedPath := "/usr/local/bin/pomodoro"
	
	if _, err := os.Stat(installedPath); err == nil {
		fmt.Printf("Found binary at %s\n", installedPath)
		fmt.Println("Note: Removing the binary requires sudo privileges.")
		fmt.Printf("Please run: sudo rm %s\n", installedPath)
	} else {
		fmt.Printf("Binary not found at %s (may already be removed)\n", installedPath)
	}
	
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error: Could not determine home directory")
		return
	}
	
	configDir := filepath.Join(homeDir, ".pomodoro")
	if _, err := os.Stat(configDir); err == nil {
		fmt.Printf("\nRemoving configuration directory: %s\n", configDir)
		if err := os.RemoveAll(configDir); err != nil {
			fmt.Printf("Error: Could not remove config directory: %v\n", err)
		} else {
			fmt.Println("Configuration removed successfully")
		}
	} else {
		fmt.Println("\nNo configuration directory found")
	}
	
	fmt.Println("\nUninstallation complete!")
	fmt.Println("Don't forget to remove the binary with: sudo rm /usr/local/bin/pomodoro")
}