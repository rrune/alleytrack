package models

import "time"

type config struct {
	Url           string `yaml:"url"`
	Enabled       bool   `yaml:"enabled"`
	Port          string `yaml:"port"`
	JwtKey        string `yaml:"jwtKey"`
	LogoText      string `yaml:"logoText"`
	AdminPassword string `yaml:"adminPassword"`
	RemovalDate   string `yaml:"removalDate"`
}

type Alleycat struct {
	WelcomeText string
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

type Participant struct {
	Number    int
	Name      string
	OutOfTown bool
	Flinta    bool
}
