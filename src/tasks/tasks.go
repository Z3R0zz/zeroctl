package tasks

import (
	"context"
	"time"
	"zeroctl/src/database"
	"zeroctl/src/handlers"

	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

type Task struct {
	Name        string
	Description string
	Schedule    string
	Handler     func(context.Context) error
	EntryID     cron.EntryID
}

type Scheduler struct {
	cronScheduler *cron.Cron
	tasks         map[string]*Task
}

var ErrTaskNotFound = database.ErrKeyNotFound

func (s *Scheduler) GetAvailableTasks() map[string]*Task {
	return s.tasks
}

func IsTaskEnabled(taskName string) (bool, error) {
	status, err := database.GetValue(taskName)
	if err != nil && err != database.ErrKeyNotFound {
		return false, err
	}

	return status != "disabled", nil
}

func (s *Scheduler) taskWrapper(task *Task) func() {
	return func() {
		enabled, err := IsTaskEnabled(task.Name)
		if err != nil {
			logrus.Errorf("Error checking if task %s is enabled: %v", task.Name, err)
			return
		}

		if !enabled {
			logrus.Debugf("Task %s is disabled, skipping execution", task.Name)
			return
		}

		logrus.Infof("Running task: %s", task.Name)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer cancel()

		err = task.Handler(ctx)
		if err != nil {
			logrus.Errorf("Error running task %s: %v", task.Name, err)
		}
	}
}

func (s *Scheduler) InitScheduler() {
	s.cronScheduler = cron.New(cron.WithSeconds())
	s.tasks = make(map[string]*Task)

	if s.cronScheduler == nil {
		logrus.Fatalf("Error creating scheduler")
		return
	}

	s.registerTasks()

	s.cronScheduler.Start()
	logrus.Infof("Scheduler started with %d tasks", len(s.tasks))
}

func (s *Scheduler) registerTasks() {
	tasks := []*Task{
		{
			Name:        "wallpaper",
			Description: "Refreshes the desktop wallpaper periodically",
			Schedule:    "0 */15 * * * *",
			Handler: func(ctx context.Context) error {
				logrus.Infof("Refreshing wallpaper")
				return handlers.RandomWallpaper()
			},
		},
		{
			Name:        "weather",
			Description: "Updates cached weather data periodically",
			Schedule:    "0 */15 * * * *",
			Handler: func(ctx context.Context) error {
				logrus.Infof("Refreshing weather data")
				return handlers.CacheWeatherData()
			},
		},
	}

	for _, task := range tasks {
		s.scheduleTask(task)
	}
}

func (s *Scheduler) scheduleTask(task *Task) {
	entryID, err := s.cronScheduler.AddFunc(task.Schedule, s.taskWrapper(task))
	if err != nil {
		logrus.Fatalf("Error scheduling task %s: %v", task.Name, err)
	}

	task.EntryID = entryID
	s.tasks[task.Name] = task

	logrus.Infof("Registered task: %s (%s)", task.Name, task.Description)
}

func (s *Scheduler) StopScheduler() {
	s.cronScheduler.Stop()
	logrus.Infof("Scheduler stopped")
}

func (s *Scheduler) GetTaskStatus(taskName string) (string, error) {
	if _, exists := s.tasks[taskName]; !exists {
		return "", ErrTaskNotFound
	}

	status, err := database.GetValue(taskName)
	if err != nil && err != database.ErrKeyNotFound {
		return "", err
	}

	if status == "disabled" {
		return "disabled", nil
	}

	return "enabled", nil
}
