//go:build darwin

package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func showASCIIArtWithAutoSize(filename string) error {
	width, height, _, err := parseASCIIArtDimensions(filename)
	if err != nil {
		width, height = 900, 900
	}
	return showASCIIArt(filename, width, height)
}

func parseASCIIArtDimensions(filename string) (width, height int, artStartLine int, err error) {
	file, err := asciiPaintings.Open(filename)
    if err != nil {
        file, err = os.Open(filename)
        if err != nil {
            return 0, 0, 0, err
        }
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
				pixelWidth := charWidth*8 + 100
				pixelHeight := charHeight*16 + 100
				
				return pixelWidth, pixelHeight, lineNum, nil
			}
		}		
		if lineNum > 3 {
			break
		}
	}
	return 900, 900, 0, nil
}

func showASCIIArt(filename string, width, height int) error {
	var artPath string
	
	// First try embedded filesystem
	embeddedFile, err := asciiPaintings.Open(filename)
	if err == nil {
		// Extract embedded file to temp location
		data, readErr := io.ReadAll(embeddedFile)
		embeddedFile.Close()
		if readErr != nil {
			return readErr
		}
		
		// Create temp file
		tmpFile, tmpErr := os.CreateTemp("", "pomodoro_ascii_*")
		if tmpErr != nil {
			return tmpErr
		}
		
		if _, writeErr := tmpFile.Write(data); writeErr != nil {
			tmpFile.Close()
			os.Remove(tmpFile.Name())
			return writeErr
		}
		tmpFile.Close()
		
		artPath = tmpFile.Name()
	} else {
		if _, err := os.Stat(filename); err == nil {
			absPath, err := filepath.Abs(filename)
			if err != nil {
				return err
			}
			artPath = absPath
		} else {
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
	}
	
	cleanupCmd := ""
	if strings.HasPrefix(artPath, os.TempDir()) {
		cleanupCmd = fmt.Sprintf("rm -f '%s'", artPath)
	}
	
	scriptContent := fmt.Sprintf(`#!/bin/bash
# Function to close window on exit
close_window() {
    %s
    # Small delay to ensure script has fully exited
    sleep 0.1
    # Close the terminal window
    osascript -e 'tell application "Terminal" to close front window' 2>/dev/null
}

# Set trap to close window on exit
trap close_window EXIT

# Skip the DIMENSIONS metadata line if present
if head -n1 '%s' | grep -q "DIMENSIONS:"; then
    tail -n +2 '%s'
else
    cat '%s'
fi
echo ''
echo 'Press Enter to close this window...'
read -r
# Exit cleanly - trap will handle window closing
exit 0
`, cleanupCmd, artPath, artPath, artPath)
	
	// Write script to temp file
	tmpScript := "/tmp/pomodoro_ascii.sh"
	if err := os.WriteFile(tmpScript, []byte(scriptContent), 0755); err != nil {
		return err
	}
	
	left, top := 100, 100
	right := left + width
	bottom := top + height
	
	// Use AppleScript to run script with custom window size
	// The script will exit and close the window via trap
	appleScript := fmt.Sprintf(`
		tell application "Terminal"
			activate
			do script "%s"
			set bounds of window 1 to {%d, %d, %d, %d}
		end tell
	`, tmpScript, left, top, right, bottom)
	
	cmd := exec.Command("osascript", "-e", appleScript)
	return cmd.Start()
}
