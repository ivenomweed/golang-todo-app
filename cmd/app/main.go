package main

import (
	"log"
	"os"
)

func main() {
	logFile, err := os.OpenFile(
		"./tmp/logs",
		os.O_APPEND|os.O_RDWR|os.O_CREATE,
		0644,
	)
	if err != nil {
		log.Panic(err)
	}

	defer func() {
		err := logFile.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	log.SetOutput(logFile)
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}
