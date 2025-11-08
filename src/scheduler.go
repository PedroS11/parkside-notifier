package main

import (
	"time"

	"github.com/go-co-op/gocron/v2"
)

func CreateCronJob(job func()) (gocron.Scheduler, error) {
	s, err := gocron.NewScheduler()
	if err != nil {
		LogError("Scheduler CreateCronJob", err)
		return nil, err
	}

	_, err = s.NewJob(
		gocron.DurationJob(
			24*time.Hour,
		),
		gocron.NewTask(
			job,
		),
	)

	if err != nil {
		LogError("Scheduler error:", err)
		return nil, err
	}

	return s, nil
}
