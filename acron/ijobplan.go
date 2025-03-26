package acron

import (
	"github.com/go-co-op/gocron/v2"
	"github.com/gofrs/uuid/v5"
)

type IJobPlanWrapper interface {
}

// use cases for this
// webui to manage cron shell/scripts (ESRI)
// * do we need to know if the job ran succesfully? errored? YES
// * log that at a app level? YES
// * do we need the ability at the app level to associate the CronJob with an app JobDefinition? YES (use GetCronJobId and GetJobDefinitionId)
// * do we need to get the last saved "state" of data? YES

type IJobPlan interface {
	GetJobPlanId() uuid.UUID
	GetTitle() string
	Validate() error
	SetupGoCronJob() (gocron.JobDefinition, []gocron.JobOption, error)

	GetRunFunction() (function any)
	Run(ccc ICronControlCenter) (IJRun, error)

	GetCronJobId() uuid.UUID
	SetCronJobId(uuid uuid.UUID)

	GetLastJRun() IJRun
	SetLastJRun(jRun IJRun)

	GetTask() ITask

	// Optional
	GetFilePath() string
	SetFilePath(filePath string)

	// These are not required here.
	//GetStashedData(key string) (interface{}, error)
	//SetStashedData(key string, data interface{}) error
}

type IJobPlans []IJobPlan

type IJobCCC interface {
	GetCCC() ICronControlCenter
}

//// runJob runs the provided ITask interface implementation.
//func runJob(job IJob) {
//	job.LogBeginJob()
//	if err := job.Task.Run(); err != nil {
//		fmt.Printf("error running job: %v\n", err)
//	} else {
//		fmt.Println("job ran successfully")
//	}
//	job.LogEndJob()
//}

//create a second file, which is the jobRunLog file
//{
//"crontab": "*/5 * * * *",
//"runImmediately": true,
//"isOneTime": false,
//"startAt": null,
//"endAt": null,
//"task": {
//"type": "shell",
//"scriptToRun": "script1.py"
//}
//}
