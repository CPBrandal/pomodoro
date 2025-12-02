package main

import "time"

const DEFAULT_WORK_DURATION = 25
const DEFAULT_BREAK_DURATION = 5
const DEFAULT_LONGER_BREAK_DURATION = 20

type LastPresetChoice struct {
	workTime time.Duration
	breakTime time.Duration
	longerBreakTime time.Duration
}

type Preset struct {
	Name         string `json:"name"`
	WorkMinutes  int    `json:"workMinutes"`
	BreakMinutes int    `json:"breakMinutes"`
	LongerBreakMinutes int    `json:"longerBreakMinutes"`
}

type Presets struct {
	Presets        []Preset `json:"presets"`
	LastUsedPreset *Preset  `json:"lastUsedPreset,omitempty"`
}