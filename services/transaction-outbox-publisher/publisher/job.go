package publisher

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog/log"
	"os"
)

type Job struct {
	scheduler *cron.Cron
	service   Service
}

func NewJob(service Service) *Job {
	return &Job{
		scheduler: cron.New(cron.WithSeconds()),
		service:   service,
	}
}

func (j *Job) Start() error {
	cronExp := os.Getenv("FREQUENCY")
	if cronExp == "" {
		return fmt.Errorf("FREQUENCY environment variable is not set")
	}

	// Add the cron job to trigger the publisher service periodically.
	_, err := j.scheduler.AddFunc(cronExp, func() {
		eventCount, err := j.service.TriggerPublisher()
		if err != nil {
			log.Error().Err(err).Msg("Cron job: Failed to trigger publisher")
		} else {
			log.Info().Int("event_count", eventCount).Msg("Cron job: Successfully triggered publisher")
		}
	})
	if err != nil {
		return err
	}

	// Start the cron scheduler
	j.scheduler.Start()
	return nil
}

// Stop stops the cron scheduler.
func (j *Job) Stop() {
	j.scheduler.Stop()
}
