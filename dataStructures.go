package main

import "time"

const DEFAULT_WORK_DURATION = 25
const DEFAULT_BREAK_DURATION = 5
const DEFAULT_LONGER_BREAK_DURATION = 20

var artworkList = []Artwork{
	{"The Persistence of Memory", "ascii-paintings/thePersistenceOfMemory", 88},
	{"Four Dancers", "ascii-paintings/fourDancers", 73},
	{"Mona Lisa", "ascii-paintings/monaLisa", 48},
	{"The Birth of Venus", "ascii-paintings/theBirthOfVenus", 70},
	{"The Death of Marat", "ascii-paintings/theDeathOfMarat", 60},
	{"Girl with a Pearl Earring", "ascii-paintings/theGirlWithThePearlEarring", 63},
	{"The Picnic", "ascii-paintings/ThePicnic", 71},
	{"The Scream", "ascii-paintings/theScream", 134},
	{"Woman with a Parasol", "ascii-paintings/womanWithAParasol", 118},
}

type LastPresetChoice struct {
	WorkTime time.Duration
	BreakTime time.Duration
	LongerBreakTime time.Duration
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

type ArtworkProgress struct {
	CurrentArtworkIndex int            `json:"currentArtworkIndex"`
	UnlockedLines       map[string]int `json:"unlockedLines"` // artwork name -> number of unlocked lines
}

type Artwork struct {
	Name     string
	Filename string
	TotalLines int
}

type UsageStats struct {
	TotalHours float64 `json:"totalHours"`
	TotalPomodoros int `json:"totalPomodoros"`
	LastUpdated time.Time `json:"lastUpdated"`
}