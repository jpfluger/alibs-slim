package acron

import (
	"github.com/jpfluger/alibs-slim/autils"
	"testing"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/stretchr/testify/assert"
	"github.com/jpfluger/alibs-slim/aerr"
)

func TestNewJRun(t *testing.T) {
	jRun := NewJRun()
	assert.NotNil(t, jRun)
	assert.NotNil(t, jRun.Logger())
	assert.Equal(t, uuid.Nil, jRun.GetJobPlanId())
	assert.Equal(t, "", jRun.GetJobPlanTitle())
}

func TestNewJRunWithOptions(t *testing.T) {
	jobPlanId := autils.NewUUID()
	jobPlanTitle := "Test Job Plan"
	jRun := NewJRunWithOptions(jobPlanId, jobPlanTitle, TASKTYPE_SHELL)
	assert.NotNil(t, jRun)
	assert.NotNil(t, jRun.Logger())
	assert.Equal(t, jobPlanId, jRun.GetJobPlanId())
	assert.Equal(t, jobPlanTitle, jRun.GetJobPlanTitle())
	assert.Equal(t, jRun.GetTaskType(), TASKTYPE_SHELL)
}

func TestBeginAndEnd(t *testing.T) {
	jRun := NewJRun()
	jRun.Begin()
	assert.False(t, jRun.IsFinished())
	assert.WithinDuration(t, time.Now(), jRun.GetStartTime(), time.Second)

	jRun.End()
	assert.True(t, jRun.IsFinished())
	assert.WithinDuration(t, time.Now().UTC(), *jRun.GetEndTime(), time.Second)
}

func TestLogging(t *testing.T) {
	jobPlanId := autils.NewUUID()
	jRun := NewJRunWithOptions(jobPlanId, "Test Job Plan", TASKTYPE_SHELL)
	logger := jRun.Logger()

	logger.Info().Msg("Test log entry")
	logs := jRun.GetLogs()
	assert.NotEmpty(t, logs)
	assert.Contains(t, logs[0], "Test log entry")
	assert.Contains(t, logs[0], jobPlanId.String())
}

func TestGetError(t *testing.T) {
	jRun := NewJRun()
	assert.Nil(t, jRun.GetError())

	err := aerr.New("Test error")
	jRun.Error = err
	assert.Equal(t, err.ToError(), jRun.GetError())
}

func TestGetLogs(t *testing.T) {
	jRun := NewJRun()
	logger := jRun.Logger()

	logger.Info().Msg("First log entry")
	logger.Info().Msg("Second log entry")

	logs := jRun.GetLogs()
	assert.Len(t, logs, 2)
	assert.Contains(t, logs[0], "First log entry")
	assert.Contains(t, logs[1], "Second log entry")
}

func TestGetByJobPlanId(t *testing.T) {
	jobPlanId := autils.NewUUID()
	runs := IJRuns{
		&JRun{JobPlanId: jobPlanId},
		&JRun{JobPlanId: autils.NewUUID()},
	}

	result := runs.GetByJobPlanId(jobPlanId)
	assert.Len(t, result, 1)
	assert.Equal(t, jobPlanId, result[0].GetJobPlanId())
}

func TestGetByJobPlanTitle(t *testing.T) {
	jobPlanTitle := "Test Plan"
	runs := IJRuns{
		&JRun{JobPlanTitle: jobPlanTitle},
		&JRun{JobPlanTitle: "Other Plan"},
	}

	result := runs.GetByJobPlanTitle(jobPlanTitle)
	assert.Len(t, result, 1)
	assert.Equal(t, jobPlanTitle, result[0].GetJobPlanTitle())
}

func TestGetByTaskType(t *testing.T) {
	taskType := TaskType("TestType")
	runs := IJRuns{
		&JRun{TaskType: taskType},
		&JRun{TaskType: TaskType("OtherType")},
	}

	result := runs.GetByTaskType(taskType)
	assert.Len(t, result, 1)
	assert.Equal(t, taskType, result[0].GetTaskType())
}

func TestGetFinished(t *testing.T) {
	runs := IJRuns{
		&JRun{},
		&JRun{},
		&JRun{},
	}
	runs[0].Begin() // Valid
	runs[0].End()
	runs[1].Begin() // Not finished.
	runs[2].End()   // Valid finished but with error.

	result := runs.GetFinished()
	assert.Len(t, result, 2)
	assert.True(t, result[0].IsFinished())
	assert.True(t, result[1].IsFinished())
}

func TestGetNotFinished(t *testing.T) {
	runs := IJRuns{
		&JRun{},
		&JRun{},
		&JRun{},
	}
	runs[0].Begin()
	runs[0].End()
	runs[1].Begin()
	runs[2].End()

	result := runs.GetNotFinished()
	assert.Len(t, result, 1)
	assert.False(t, result[0].IsFinished())
}

func TestReplaceByJobPlanId(t *testing.T) {
	jobPlanId := autils.NewUUID()
	runs := IJRuns{
		&JRun{JobPlanId: jobPlanId, JobPlanTitle: "Old Title", TaskType: TASKTYPE_SHELL},
	}

	newRun := &JRun{JobPlanId: jobPlanId, JobPlanTitle: "New Title", TaskType: "new-type"}
	runs.ReplaceByJobPlanId(newRun)

	assert.Len(t, runs, 1)
	assert.Equal(t, "New Title", runs[0].GetJobPlanTitle())
	assert.Equal(t, "new-type", runs[0].GetTaskType().String())
}
