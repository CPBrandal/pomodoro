//go:build !darwin && !linux

package main

import (
	"bufio"
	"fmt"
	"os"
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
	// Fallback: just print to current terminal (size parameters ignored)
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	fmt.Print(string(data))
	return nil
}
