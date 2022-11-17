package main

import (
	"context"
	"encoding/json"
	"github.com/pkg/browser"
	"github.com/procyon-projects/chrono"
	"io"
	"log"
	"os"
	"os/signal"
	"time"
)

type Config struct {
	Time  []string `json:"time"`
	Links []string `json:"links"`
}

func joinMeet(index int, link string) func(ctx context.Context) {
	return func(ctx context.Context) {
		err := browser.OpenURL(link)
		if err != nil {
			panic(err)
		} else {
			log.Println("Joined meeting", index+1)
		}
	}
}

func main() {
	var config *Config

	now := time.Now()

	scheduler := chrono.NewDefaultTaskScheduler()

	file, err := os.Open("./config.json")
	if err != nil {
		panic(err)
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	fileContent, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(fileContent, &config)
	if err != nil {
		panic(err)
	}

	for idx, configTime := range config.Time {
		parsedTime, err := time.ParseInLocation(time.Kitchen, configTime, time.Local)
		parsedTime = parsedTime.Add(-time.Minute * 5)
		startTime := time.Date(now.Year(), now.Month(), now.Day(), parsedTime.Hour(), parsedTime.Minute(), parsedTime.Second(), parsedTime.Nanosecond(), time.Local)
		if err != nil {
			panic(err)
		}
		_, err = scheduler.Schedule(joinMeet(idx, config.Links[idx]), chrono.WithTime(startTime))
		if err != nil {
			panic(err)
		} else {
			log.Println("Added schedule", idx+1, "on", startTime.Format(time.Kitchen))
		}
	}

	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, os.Kill)
	<-ch
}
