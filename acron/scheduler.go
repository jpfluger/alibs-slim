package acron

import (
	"fmt"
	"path/filepath"
	"reflect"
	"sync"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/jpfluger/alibs-slim/ajson"

	gocron "github.com/go-co-op/gocron/v2"
)

// globalCron is a globally accessible gocron.Scheduler instance.
var (
	globalCron gocron.Scheduler
	once       sync.Once
	mutex      sync.Mutex // Mutex to protect access to globalCron
)

// SCHEDULER returns the global gocron.Scheduler instance.
// It ensures that the scheduler is initialized with UTC location.
func SCHEDULER() gocron.Scheduler {
	once.Do(func() {
		// Initialize the scheduler with UTC location only once.
		globalCron, _ = gocron.NewScheduler(gocron.WithLocation(time.UTC))
	})
	return globalCron
}

// SetScheduler safely sets the global scheduler instance.
func SetScheduler(scheduler gocron.Scheduler, doReinitWithShutdown bool) error {
	mutex.Lock()         // Lock the mutex before modifying globalCron.
	defer mutex.Unlock() // Unlock the mutex when the function returns.

	if doReinitWithShutdown {
		if err := SCHEDULER().Shutdown(); err != nil {
			return fmt.Errorf("failed shutdown scheduler: %w", err)
		}
	}
	if scheduler == nil {
		scheduler, _ = gocron.NewScheduler(gocron.WithLocation(time.UTC))
	}
	globalCron = scheduler
	return nil
}

// FindJobJSONFiles finds 'job.json' files in subdirectories of the given directory.
// It searches subdirectories only one level deep.
func FindJobJSONFiles(workingDir string) ([]string, error) {
	pattern := filepath.Join(workingDir, "*", "job.json") // Pattern to match 'job.json' files.
	return filepath.Glob(pattern)                         // Find files matching the pattern.
}

// LoadJobJSONFiles loads job definitions from 'job.json' files found in the working directory.
func LoadJobJSONFiles(workingDir string, jobPlanType reflect.Type) (IJobPlans, error) {
	files, err := FindJobJSONFiles(workingDir)
	if err != nil {
		return nil, fmt.Errorf("error finding job json files: %v", err)
	}
	if len(files) == 0 {
		return nil, nil // No files found, return nil.
	}

	var jobs IJobPlans
	for _, file := range files {
		job := reflect.New(jobPlanType).Interface().(IJobPlan)
		if err := ajson.UnmarshalFile(file, job); err != nil {
			return nil, fmt.Errorf("error unmarshaling job json file: %v", err)
		}
		if err := job.Validate(); err != nil {
			return nil, fmt.Errorf("error validating job: %v", err)
		}
		job.SetFilePath(file)
		jobs = append(jobs, job)
	}
	return jobs, nil
}

// AddJobPlan adds a job to the globalCron scheduler.
func AddJobPlan(jobPlan IJobPlan) error {
	jobDef, jobOptions, err := jobPlan.SetupGoCronJob()
	if err != nil {
		return fmt.Errorf("error setting up gocron job: %v", err)
	}

	mutex.Lock()         // Lock the mutex before adding the job.
	defer mutex.Unlock() // Unlock the mutex when the function returns.

	fn := jobPlan.GetRunFunction()
	if fn == nil {
		fn = runJobPlanDefault
	}

	// Add the job to the scheduler with the defined options.
	//_, err = globalCron.NewJob(jobDef, gocron.NewTask(runITask, job.Task), jobOptions...)
	job, err := globalCron.NewJob(jobDef, gocron.NewTask(fn, jobPlan), jobOptions...)
	if err != nil {
		return fmt.Errorf("error adding job to globalCron: %v", err)
	}
	// Convert to "gofrs/uuid/v5"
	jobPlan.SetCronJobId(uuid.UUID(job.ID()))

	return nil
}

// ScheduleJobPlans schedules a list of jobs using the globalCron scheduler.
func ScheduleJobPlans(jobPlans IJobPlans) error {
	for _, jobPlan := range jobPlans {
		if err := AddJobPlan(jobPlan); err != nil {
			return fmt.Errorf("error scheduling job: %v", err)
		}
	}
	return nil
}

// runITask runs the provided ITask interface implementation.
func runITask(job ITask) {
	if err := job.Run(&CronControlCenter{}); err != nil {
		fmt.Printf("error running job: %v\n", err)
	} else {
		fmt.Println("job ran successfully")
	}
}

// runJobPlanDefault runs the provided IJobPlan interface implementation.
// Copy this to your app level and compose your own version of CronControlCenter
// for fully-customized control over running a job from the app level. It is
// here that the app can retrieve or save the state of last JRuns, can associate
// specific databases or other API/Modules with a task and can handle the results
// of the task, if the App requires it.
// See tests  for CronControlCenterShell as a working blueprint.
func runJobPlanDefault(jobPlan IJobPlan) {
	ccc := &CronControlCenter{
		jrun: NewJRunWithOptions(jobPlan.GetJobPlanId(), jobPlan.GetTitle(), jobPlan.GetTask().GetType()),
	}
	_, _ = jobPlan.Run(ccc)
}
