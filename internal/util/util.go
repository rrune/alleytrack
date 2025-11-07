package util

import (
	"log"
	"os"

	"github.com/rrune/alleytrack/pkg/models"
	"gopkg.in/yaml.v2"
)

func Check(err error) (r bool) {
	if err != nil {
		r = true
	}
	return
}

func CheckWLogs(err error) (r bool) {
	if err != nil {
		r = true
		log.Println(err)
	}
	return
}

func CheckPanic(err error) {
	if err != nil {
		f, err2 := os.OpenFile("./data/err.txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
		if err2 != nil {
			log.Fatal(err, err2)
		}
		defer f.Close()
		_, err2 = f.Write([]byte(err.Error()))
		if err2 != nil {
			log.Fatal(err, err2)
		}
		panic(err)
	}
}

func SwitchEnabledInConfig(alleycat *models.Alleycat) (err error) {
	alleycat.Config.Enabled = !alleycat.Config.Enabled

	confByte, err := yaml.Marshal(alleycat.Config)
	if Check(err) {
		return
	}
	err = os.WriteFile("./config/config.yml", confByte, 0644)
	return
}

func WriteEvent(event string) (err error) {
	f, err := os.OpenFile("./data/event.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if Check(err) {
		return
	}
	defer f.Close()
	_, err = f.WriteString("\n" + event)
	return
}
