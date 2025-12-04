//go:build linux

package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func showASCIIArtWithAutoSize(filename string) error {
	width, height, _, err := parseASCIIArtDimensions(filename)
	if err != nil {
		// Use defaults if parsing fails
		width, height = 900, 900
	}
	return showASCIIArt(filename, width, height)
}

func parseASCIIArtDimensions(filename string) (width, height int, artStartLine int, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, 0, 0, err
	}
	defer file.Close()
	
	scanner := bufio.NewScanner(file)
	lineNum := 0
	
	for scanner.Scan() {
		line := scanner.Text()
		lineNum++
		
		// Check for DIMENSIONS metadata
		if strings.HasPrefix(line, "DIMENSIONS:") {
			dimStr := strings.TrimPrefix(line, "DIMENSIONS:")
			parts := strings.Split(dimStr, "x")
			if len(parts) == 2 {
				charWidth, _ := strconv.Atoi(strings.TrimSpace(parts[0]))
				charHeight, _ := strconv.Atoi(strings.TrimSpace(parts[1]))
				
				// Convert character dimensions to pixel dimensions
				// Terminal character size: ~8px wide x ~16px tall
				// Add window chrome: ~100px for borders/title bar
				pixelWidth := charWidth*8 + 100   // characters * 8px + window chrome
				pixelHeight := charHeight*16 + 100 // lines * 16px + window chrome
				
				return pixelWidth, pixelHeight, lineNum, nil
			}
		}
		
		// Stop after first few lines if no metadata found
		if lineNum > 3 {
			break
		}
	}
	
	// Default dimensions if not found
	return 900, 900, 0, nil
}

func showASCIIArt(filename string, width, height int) error {
	// Get the absolute path to the ASCII art file
	exePath, err := os.Executable()
	if err != nil {
		return err
	}
	exeDir := filepath.Dir(exePath)
	artPath := filepath.Join(exeDir, filename)
	
	// Check if file exists
	if _, err := os.Stat(artPath); err != nil {
		return err
	}
	
	// Try different terminal emulators with geometry support
	terminals := []string{"gnome-terminal", "xterm", "konsole", "x-terminal-emulator"}
	
	for _, term := range terminals {
		if _, err := exec.LookPath(term); err == nil {
			var cmd *exec.Cmd
			geometry := fmt.Sprintf("%dx%d", width, height)
			switch term {
			case "gnome-terminal":
				cmd = exec.Command(term, "--geometry", geometry, "--", "bash", "-c", fmt.Sprintf("cat '%s'; echo ''; read -p 'Press Enter to close...';", artPath))
			case "xterm":
				cmd = exec.Command(term, "-geometry", geometry, "-e", "bash", "-c", fmt.Sprintf("cat '%s'; echo ''; read -p 'Press Enter to close...';", artPath))
			case "konsole":
				cmd = exec.Command(term, "--geometry", geometry, "-e", "bash", "-c", fmt.Sprintf("cat '%s'; echo ''; read -p 'Press Enter to close...';", artPath))
			default:
				cmd = exec.Command(term, "-e", "bash", "-c", fmt.Sprintf("cat '%s'; echo ''; read -p 'Press Enter to close...';", artPath))
			}
			return cmd.Start() // Start in background
		}
	}
	
	return fmt.Errorf("no terminal emulator found")
}