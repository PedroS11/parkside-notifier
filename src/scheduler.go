package main

import (
	"log/slog"
	"time"

	"github.com/go-co-op/gocron/v2"
)

func CreateCronJob(job func()) (gocron.Scheduler, error) {
	s, err := gocron.NewScheduler()
	if err != nil {
		slog.Error("Scheduler CreateCronJob", "error", err.Error())
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
		slog.Error("Scheduler NewJob", "error", err.Error())
		return nil, err
	}

	return s, nil
}
