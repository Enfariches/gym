package gcron

import (
	"fmt"
	"health/internal/domain/models"
	"sync"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
)

type JobScheduler struct {
	scheduler gocron.Scheduler
	jobs      map[int64]uuid.UUID
	mu        sync.Mutex
	stopCh    chan struct{}
}

func NewJobScheduler() (*JobScheduler, error) {
	jobScheduler, err := gocron.NewScheduler()
	if err != nil {
		return nil, fmt.Errorf("error creating Scheduler: %v", err)
	}

	return &JobScheduler{
		scheduler: jobScheduler,
		jobs:      make(map[int64]uuid.UUID),
		stopCh:    make(chan struct{}),
	}, nil
}

func (j *JobScheduler) Start() {
	j.scheduler.Start()
	go j.run()
}

func (j *JobScheduler) run() {
	<-j.stopCh
}

func (j *JobScheduler) Stop() {
	close(j.stopCh)
	j.scheduler.StopJobs()
}

// fmt.Printf's пока нужны для разработки и контроля работы планировщика
func (j *JobScheduler) UpdateScheduler(schedules []*models.Schedule) {
	j.mu.Lock()
	defer j.mu.Unlock()
	for _, v := range schedules {
		if jobId, ok := j.jobs[v.ID]; ok {
			j.scheduler.RemoveJob(jobId)
			delete(j.jobs, v.ID)

			fmt.Printf("Job with ID:%v removed\n", jobId)
		}

		if !v.IsActive {
			continue
		}

		activeJob, err := j.scheduler.NewJob(gocron.CronJob(v.CronExpression, true), gocron.NewTask(work, "hello", v.MediaID))
		if err != nil {
			fmt.Println("Error creating job:", err)
		}

		fmt.Printf("Job with ID:%v added \n", activeJob.ID())
		j.jobs[v.ID] = activeJob.ID()

	}
}

func work(a string, MediaID int64) {
	fmt.Printf("Running task: %s with argument %d\n", a, MediaID)
}
