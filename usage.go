package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type UsageStats struct {
	TotalHours float64 `json:"totalHours"`
	LastUpdated time.Time `json:"lastUpdated"`
}

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
		return UsageStats{TotalHours: 0}
	}

	data, err := os.ReadFile(usagePath)
	if err != nil {
		return UsageStats{TotalHours: 0}
	}

	var stats UsageStats
	if err := json.Unmarshal(data, &stats); err != nil {
		return UsageStats{TotalHours: 0}
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

func addSessionTime(duration time.Duration) error {
	stats := loadUsageStats()
	stats.TotalHours += duration.Hours()
	return saveUsageStats(stats)
}

func getTotalUsageHours() float64 {
	stats := loadUsageStats()
	return stats.TotalHours
}
