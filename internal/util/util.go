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
		f, err2 := os.OpenFile("./data/err.txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
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
