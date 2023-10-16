package main

import (
	"fmt"
	"io"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{
		TimestampFormat: "2006-01-02T15:04:05",
	})
	log.SetLevel(log.DebugLevel)

	if _, err := os.Stat("./logs"); os.IsNotExist(err) {
		if err := os.MkdirAll("./logs", 0666); err != nil {
			log.Warn(err)
		}
	}

	if file, err := os.OpenFile(
		fmt.Sprintf("./logs/%s.log", time.Now().Local().Format("2006-01")), 
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 
		0666,
	); err == nil {
		writers := []io.Writer{
			file,
			os.Stdout,
		}
		log.SetOutput(io.MultiWriter(writers...))
	} else {
		log.Warn("Failed to log to file, using default stderr")
	}

	
}

func main() {
	log.Info("hello world")
	log.Debug("I love you Nahida")
}
