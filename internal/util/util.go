package util

import (
	"log"
	"os"
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

func WriteEvent(event string) (err error) {
	f, err := os.OpenFile("./data/event.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if Check(err) {
		return
	}
	defer f.Close()
	_, err = f.WriteString("\n" + event)
	return
}
