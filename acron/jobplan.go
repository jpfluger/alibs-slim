package acron

import (
	"encoding/json"
	"fmt"
	"github.com/go-co-op/gocron/v2"
	"github.com/gofrs/uuid/v5"
	"github.com/robfig/cron/v3"
	"github.com/jpfluger/alibs-slim/areflect"
	"github.com/jpfluger/alibs-slim/autils"
	"reflect"
	"strings"
	"sync"
	"time"
)

// JobPlan represents a job definition for a Cron Scheduler with options and data.
type JobPlan struct {
	Crontab        string     `json:"crontab,omitempty"`
	RunImmediately bool       `json:"runImmediately,omitempty"`
	RunLimit       int        `json:"runLimit,omitempty"`
	TimeZone       string     `json:"timeZone,omitempty"`
	StartAt        *time.Time `json:"startAt,omitempty"`
	EndAt          *time.Time `json:"endAt,omitempty"`
	Task           ITask      `json:"task,omitempty"`

	JobPlanId uuid.UUID `json:"jobPlanId,omitempty"`
	Title     string    `json:"title,omitempty"`

	isValidated bool
	timeZone    string
	utcStartAt  time.Time
	utcEndAt    time.Time

	cronJobId uuid.UUID

	// optional, if job.json is loaded from disk, then save it here.
	filePath string

	//ccc      ICronControlCenter
	lastJRun IJRun

	mu sync.RWMutex // Mutex for concurrency safety
}

// Validate the job plan.
func (j *JobPlan) Validate() error {
	j.mu.Lock()
	defer j.mu.Unlock()
	return j.validate()
}

func (j *JobPlan) validate() error {
	// Force JobPlanId to a uuid rather than error.
	if j.JobPlanId == uuid.Nil {
		j.JobPlanId = autils.NewUUID()
	}
	j.Title = strings.TrimSpace(j.Title)

	// Trim spaces from the Crontab field.
	j.Crontab = strings.TrimSpace(j.Crontab)
	if j.Crontab != "" {
		// Validate the crontab expression using the cron package.
		if _, err := cron.ParseStandard(j.Crontab); err != nil {
			return fmt.Errorf("invalid crontab expression: %v", err)
		}
	}

	// Ensure RunLimit is non-negative.
	if j.RunLimit < 0 {
		j.RunLimit = 0
	}

	// Convert StartAt and EndAt to UTC.
	var err error
	j.timeZone = strings.TrimSpace(j.TimeZone)
	if j.timeZone == "" {
		j.timeZone = "UTC"
	}
	loc, err := time.LoadLocation(j.timeZone)
	if err != nil {
		return fmt.Errorf("invalid timezone: %v", err)
	}

	if j.StartAt != nil {
		j.utcStartAt = j.StartAt.In(loc).UTC()
		if j.utcStartAt.Before(time.Now().UTC()) {
			return fmt.Errorf("startAt must not be in the past")
		}
	}

	if j.EndAt != nil {
		j.utcEndAt = j.EndAt.In(loc).UTC()
		if j.utcEndAt.Before(time.Now().UTC()) {
			return fmt.Errorf("endAt must not be in the past")
		}
		if !j.utcStartAt.IsZero() && j.utcStartAt.After(j.utcEndAt) {
			return fmt.Errorf("startAt must be before endAt")
		}
	}

	if j.Crontab == "" && j.utcStartAt.IsZero() && !j.RunImmediately {
		if j.RunLimit == 1 {
			j.RunImmediately = true
		} else {
			return fmt.Errorf("crontab, startAt, or runImmediately is required")
		}
	}

	// Validate the Task field.
	if j.Task == nil {
		return fmt.Errorf("task is required")
	}
	if err := j.Task.Validate(); err != nil {
		return fmt.Errorf("failed task validation: %v", err)
	}

	j.isValidated = true

	return nil
}

// SetupGoCronJob sets up the job with gocron using the provided options.
func (j *JobPlan) SetupGoCronJob() (gocron.JobDefinition, []gocron.JobOption, error) {
	j.mu.Lock()
	defer j.mu.Unlock()

	if !j.isValidated {
		if err := j.validate(); err != nil {
			return nil, nil, fmt.Errorf("failed validation: %v", err)
		}
	}

	var options []gocron.JobOption

	// WithSingletonMode keeps the job from running again if it is already running.
	// This is useful for jobs that should not overlap, and that occasionally
	// (but not consistently) run longer than the interval between job runs.
	//
	// LimitModeReschedule causes jobs reaching the limit set in
	// WithLimitConcurrentJobs or WithSingletonMode to be skipped
	// and rescheduled for the next run time rather than being
	// queued up to waITaskit.
	options = append(options, gocron.WithSingletonMode(gocron.LimitModeReschedule))

	// WithLimitedRuns limits the number of executions of this job to n.
	// Upon reaching the limit, the job is removed from the scheduler.
	if j.RunLimit > 0 {
		options = append(options, gocron.WithLimitedRuns(uint(j.RunLimit)))
	}

	// WithStartImmediately tells the scheduler to run the job immediately
	// regardless of the type or schedule of job. After this immediate run
	// the job is scheduled from this time based on the job definition.
	if j.RunImmediately {
		options = append(options, gocron.JobOption(gocron.WithStartImmediately()))
	}

	// WithStartDateTime sets the first date & time at which the job should run.
	// This datetime must be in the future.
	if !j.utcStartAt.IsZero() {
		options = append(options, gocron.JobOption(gocron.WithStartDateTime(j.utcStartAt)))
	}

	// WithStopDateTime sets the final date & time after which the job should stop.
	// This must be in the future and should be after the startTime (if specified).
	// The job's final run may be at the stop time, but not after.
	if !j.utcEndAt.IsZero() {
		options = append(options, gocron.JobOption(gocron.WithStopDateTime(j.utcEndAt)))
	}

	if len(options) == 0 && j.Crontab == "" {
		return nil, nil, fmt.Errorf("no valid job options")
	}

	// CronJob defines a new job using the crontab syntax: `* * * * *`.
	// Uses "github.com/robfig/cron/v3".
	// An optional 6th field can be used at the beginning if withSeconds
	// is set to true: `* * * * * *`.
	// The timezone can be set on the Scheduler using WithLocation, or in the
	// crontab in the form `TZ=America/Chicago * * * * *` or
	// `CRON_TZ=America/Chicago * * * * *`
	var jobDef gocron.JobDefinition
	if j.Crontab != "" {
		jobDef = gocron.CronJob(j.Crontab, false)
	} else {
		// OneTimeJob is to run a job once at a specified time and not on
		// any regular schedule.
		if j.utcStartAt.IsZero() {
			jobDef = gocron.OneTimeJob(gocron.OneTimeJobStartImmediately())
		} else {
			jobDef = gocron.OneTimeJob(gocron.OneTimeJobStartDateTime(j.utcStartAt))
		}
	}

	return jobDef, options, nil
}

func (j *JobPlan) GetJobPlanId() uuid.UUID {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return j.JobPlanId
}

func (j *JobPlan) GetTitle() string {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return j.Title
}

func (j *JobPlan) GetRunFunction() (function any) {
	return nil
}

func (j *JobPlan) GetCronJobId() uuid.UUID {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return j.cronJobId
}

func (j *JobPlan) SetCronJobId(cronJobId uuid.UUID) {
	j.mu.Lock()
	defer j.mu.Unlock()
	j.cronJobId = cronJobId
}

func (j *JobPlan) GetLastJRun() IJRun {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return j.lastJRun
}

func (j *JobPlan) SetLastJRun(jRun IJRun) {
	j.mu.Lock()
	defer j.mu.Unlock()
	j.lastJRun = jRun
}

func (j *JobPlan) GetTask() ITask {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return j.Task
}

func (j *JobPlan) GetFilePath() string {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return j.filePath
}

func (j *JobPlan) SetFilePath(filePath string) {
	j.mu.Lock()
	defer j.mu.Unlock()
	j.filePath = filePath
}

func (j *JobPlan) RunJobPlanDefault(ccc ICronControlCenter) (IJRun, error) {

	// Checks
	if ccc == nil {
		return nil, fmt.Errorf("ccc is nil")
	}
	if ccc.GetJRun() == nil {
		return nil, fmt.Errorf("jrun not found in job plan '%s'", j.Title)
	}

	// Begin logging
	ccc.GetJRun().Begin()

	// Run
	err := j.GetTask().Run(ccc)
	if err != nil {
		if ccc.GetJRun().GetError() == nil {
			ccc.GetJRun().SetError(err)
		}
	}

	// End logging
	ccc.GetJRun().End()

	j.mu.Lock()
	j.lastJRun = ccc.GetJRun()
	j.mu.Unlock()

	return ccc.GetJRun(), err
}

// UnmarshalJSONTask is a custom unmarshaller for JobPlan that handles ITask.
func (j *JobPlan) UnmarshalJSONTask(task json.RawMessage) error {
	j.mu.Lock()
	defer j.mu.Unlock()

	var rawmap map[string]interface{}
	if err := json.Unmarshal(task, &rawmap); err != nil {
		return fmt.Errorf("failed to unmarshal ITask: %v", err)
	}

	rawType, ok := rawmap["type"].(string)
	if !ok {
		return fmt.Errorf("type field not found or is not a string in ITask")
	}

	rtype, err := areflect.TypeManager().FindReflectType(TYPEMANAGER_CRONTASKDATA, rawType)
	if err != nil {
		return fmt.Errorf("cannot find type struct '%s': %v", rawType, err)
	}

	obj := reflect.New(rtype).Interface()
	if err = json.Unmarshal(task, obj); err != nil {
		return fmt.Errorf("failed to unmarshal ITask where type is '%s': %v", rawType, err)
	}

	iTask, ok := obj.(ITask)
	if !ok {
		return fmt.Errorf("created object does not implement ITask where type is '%s'", rawType)
	}
	j.Task = iTask

	return nil
}

// JobPlans is a slice of pointers to JobPlan structs.
type JobPlans []*JobPlan
