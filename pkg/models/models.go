package models

import "time"

type Config struct {
	Url           string `yaml:"url"`
	Port          string `yaml:"port"`
	JwtKey        string `yaml:"jwtKey"`
	LogoText      string `yaml:"logoText"`
	AdminPassword string `yaml:"adminPassword"`
	WelcomeText   string
	Manifest      []Checkpoint
}

type Checkpoint struct {
	ID           int    `yaml:"id"`
	Link         string `yaml:"link"`
	Location     string `yaml:"location"`
	Info         string `yaml:"info"`
	Text         bool   `yaml:"text"`
	Requirements []int  `yaml:"requirements"`
	Content      string
	Time         string
}

type ParticipantCheckpoint struct {
	Time    time.Time `json:"Time"`
	Content string    `json:"Content"`
}

type Participant struct {
	Number      int
	Name        string
	OutOfTown   bool
	Flinta      bool
	Checkpoints map[int]ParticipantCheckpoint
}
