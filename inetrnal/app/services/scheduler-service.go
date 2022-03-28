package services

import (
	"github.com/jasonlvhit/gocron"
	"github.com/sirupsen/logrus"
	"project/inetrnal/app/store"
	"time"
)

type SchedulerServiceI interface {
	Start()
	DoJob(job func() error, jobName string)
}

type SchedulerService struct {
	domain string
	store  *store.Store
	logger *logrus.Logger
}

func NewSchedulerService(store *store.Store, logger *logrus.Logger, domain string) SchedulerServiceI {
	return &SchedulerService{store: store, logger: logger, domain: domain}
}

func (s *SchedulerService) Start() {
	gocron.ChangeLoc(time.Local)
	cron := gocron.NewScheduler()

	time.Sleep(1 * time.Second)

	//ImportKpiDataFromDWHTable
	PrintNowJob := cron.Every(5).Second()
	PrintNowJob.Do(s.DoJob, s.PrintNow, PrintNowFunction)
	s.logger.Infof("ImportKpiDataFromDWHTable First start time: %s", PrintNowJob.NextScheduledTime().Format("02.01.2006 15:04:05"))

	<-cron.Start()
}

func (s *SchedulerService) DoJob(job func() error, jobName string) {
	startTime := time.Now()
	defer func() {
		finishTime := time.Now()
		s.logger.Infof("Finishing job: %s, finishTime: %s, duration: %s", jobName, finishTime.Format("02.01.2006 15:04:05"), finishTime.Sub(startTime))
	}()
	s.logger.Infof("Starting job: %s, startTime: %s", jobName, startTime.Format("02.01.2006 15:04:05"))

	err := job()
	if err != nil {
		s.logger.Errorf("Job %s finished with error: %v", jobName, err)
	}
}

const PrintNowFunction = "PrintNow"

func (s *SchedulerService) PrintNow() error {
	s.logger.Infof("NOW")
	return nil
}
