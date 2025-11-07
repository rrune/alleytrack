package models

import "time"

type config struct {
	Enabled       bool   `yaml:"enabled"`
	Port          string `yaml:"port"`
	JwtKey        string `yaml:"jwtKey"`
	LogoText      string `yaml:"logoText"`
	AdminPassword string `yaml:"adminPassword"`
}

type Alleycat struct {
	WelcomeText string
	Manifest    []Checkpoint
	Config      config
}

type Checkpoint struct {
	ID           int    `yaml:"id"`
	Link         string `yaml:"link"`
	Location     string `yaml:"location"`
	Info         string `yaml:"info"`
	Text         bool   `yaml:"text"`
	Requirements []int  `yaml:"requirements"`
	Content      string
	Time         time.Time
}

type ParticipantCheckpoint struct {
	Time    time.Time `json:"Time"`
	Content string    `json:"Content"`
	ID      int
}

type Participant struct {
	Number      int
	Name        string
	OutOfTown   bool
	Flinta      bool
	Checkpoints map[int]ParticipantCheckpoint
}
