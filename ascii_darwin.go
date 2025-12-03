//go:build darwin

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
	var artPath string
	
	// First try relative to current working directory (for development with go run)
	if _, err := os.Stat(filename); err == nil {
		absPath, err := filepath.Abs(filename)
		if err != nil {
			return err
		}
		artPath = absPath
	} else {
		// Try relative to executable (for installed version)
		exePath, err := os.Executable()
		if err != nil {
			return err
		}
		exeDir := filepath.Dir(exePath)
		artPath = filepath.Join(exeDir, filename)
		
		if _, err := os.Stat(artPath); err != nil {
			return fmt.Errorf("ASCII art file not found")
		}
	}
	
	// Create a script that waits for input then exits
	scriptContent := fmt.Sprintf(`#!/bin/bash
# Skip the DIMENSIONS metadata line if present
if head -n1 '%s' | grep -q "DIMENSIONS:"; then
    tail -n +2 '%s'
else
    cat '%s'
fi
echo ''
echo 'Press Enter to close this window...'
read -r
`, artPath, artPath, artPath)
	
	// Write script to temp file
	tmpScript := "/tmp/pomodoro_ascii.sh"
	if err := os.WriteFile(tmpScript, []byte(scriptContent), 0755); err != nil {
		return err
	}
	
	// Calculate window bounds: {left, top, right, bottom}
	// Default position: 100, 100
	left, top := 100, 100
	right := left + width
	bottom := top + height
	
	// Use AppleScript to run script with custom window size
	appleScript := fmt.Sprintf(`
		tell application "Terminal"
			activate
			set newTab to do script "%s"
			set bounds of window 1 to {%d, %d, %d, %d}
			repeat
				delay 0.3
				try
					if not busy of newTab then
						delay 0.5
						close newTab
						exit repeat
					end if
				on error
					exit repeat
				end try
			end repeat
		end tell
	`, tmpScript, left, top, right, bottom)
	
	cmd := exec.Command("osascript", "-e", appleScript)
	return cmd.Start()
}
