package main

import (
	"log"
	"os"

	"github.com/rrune/alleytrack/internal/util"
	"github.com/rrune/alleytrack/pkg/config"
	"github.com/rrune/alleytrack/pkg/data"
	"github.com/rrune/alleytrack/pkg/models"
	routing "github.com/rrune/alleytrack/pkg/routing"
	"gopkg.in/yaml.v2"
)

func main() {
	// create or open log file
	f, err := os.OpenFile("./data/alleytrack.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Printf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)
	log.SetFlags(2 | 3)

	// read config
	var alleycat models.Alleycat
	ymlData, err := os.ReadFile("./config/config.yml")
	util.CheckPanic(err)
	err = yaml.Unmarshal(ymlData, &alleycat.Config)
	util.CheckPanic(err)

	// read welcome text
	wel, err := os.ReadFile("./config/welcome.txt")
	util.CheckPanic(err)
	alleycat.WelcomeText = string(wel)

	// read manifest
	/* 	manifest, err := os.ReadFile("./config/manifest.yml")
	   	util.CheckPanic(err)
	   	err = yaml.Unmarshal(manifest, &alleycat.Manifest)
	   	util.CheckPanic(err) */

	// init database
	db, err := data.Init("./data/db.sqlite")
	util.CheckPanic(err)

	err = config.ReadManifestIntoDatabase("./config/manifest.yml", &db)
	util.CheckPanic(err)

	util.CheckPanic(routing.Init(&alleycat, &db))
}
