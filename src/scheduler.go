package main

import (
	"fmt"
	"time"

	"github.com/go-co-op/gocron/v2"
)

func CreateCronJob(job func()) (gocron.Scheduler, error) {
	s, err := gocron.NewScheduler()
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	_, err = s.NewJob(
		gocron.DurationJob(
			1*time.Minute,
		),
		gocron.NewTask(
			job,
		),
	)

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return s, nil
}
