package tasks

import (
	"context"
	"time"
	"zeroctl/src/handlers"

	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

type Scheduler struct {
	cronScheduler *cron.Cron
}

func (s *Scheduler) InitScheduler() {
	s.cronScheduler = cron.New(cron.WithSeconds())
	if s.cronScheduler == nil {
		logrus.Fatalf("Error creating scheduler")
		return
	}

	s.scheduleJob("0 */15 * * * *", func(ctx context.Context) error {
		logrus.Infof("Refreshing wallpaper")
		return handlers.RandomWallpaper()
	})

	s.scheduleJob("0 */15 * * * *", func(ctx context.Context) error {
		logrus.Infof("Refreshing weather data")
		return handlers.CacheWeatherData()
	})

	s.cronScheduler.Start()
	logrus.Infof("Scheduler started")
}

func (s *Scheduler) scheduleJob(cronExpr string, job func(context.Context) error) {
	_, err := s.cronScheduler.AddFunc(cronExpr, func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer cancel()

		err := job(ctx)
		if err != nil {
			logrus.Fatalf("Error running job: %v", err)
		}
	})
	if err != nil {
		logrus.Fatalf("Error scheduling job: %v", err)
	}
}

func (s *Scheduler) StopScheduler() {
	s.cronScheduler.Stop()
	logrus.Infof("Scheduler stopped")
}
