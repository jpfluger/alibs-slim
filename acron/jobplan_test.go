package acron

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestJob_Validate tests the Validate method of the Job struct.
func TestJob_Validate(t *testing.T) {
	now := time.Now().UTC()
	past := now.Add(-time.Hour)
	future := now.Add(time.Hour)

	tests := []struct {
		name    string
		job     JobPlan
		wantErr bool
	}{
		{
			name: "Valid job with future start and end dates",
			job: JobPlan{
				Crontab:        "* * * * *",
				RunImmediately: true,
				RunLimit:       1,
				StartAt:        &future,
				EndAt:          &future,
				Task: &MockITask{
					Executed: false,
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid job with past start date",
			job: JobPlan{
				StartAt: &past,
				Task:    &MockITask{Executed: false},
			},
			wantErr: true,
		},
		{
			name: "Invalid job with start date after end date",
			job: JobPlan{
				StartAt: &future,
				EndAt:   &past,
				Task:    &MockITask{Executed: false},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.job.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestJob_ValidateTimes(t *testing.T) {
	// Mock task to satisfy the Job struct requirements
	mockTask := &MockITask{}

	// Define a fixed point in time for testing
	testTimeStart := time.Date(2030, 3, 29, 15, 4, 5, 0, time.UTC)
	testTimeEnd := time.Date(2030, 3, 30, 15, 4, 5, 0, time.UTC)

	// Test cases
	tests := []struct {
		name      string
		timeZone  string
		startAt   *time.Time
		endAt     *time.Time
		expectErr bool
	}{
		{
			name:      "Empty time zone defaults to UTC",
			timeZone:  "",
			startAt:   &testTimeStart,
			endAt:     &testTimeEnd,
			expectErr: false,
		},
		{
			name:      "Explicit UTC time zone",
			timeZone:  "UTC",
			startAt:   &testTimeStart,
			endAt:     &testTimeEnd,
			expectErr: false,
		},
		{
			name:      "Non-UTC time zone",
			timeZone:  "America/Chicago",
			startAt:   &testTimeStart,
			endAt:     &testTimeEnd,
			expectErr: false,
		},
		{
			name:      "Invalid time zone",
			timeZone:  "Invalid/TimeZone",
			startAt:   &testTimeStart,
			endAt:     &testTimeEnd,
			expectErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			job := JobPlan{
				TimeZone: tc.timeZone,
				StartAt:  tc.startAt,
				EndAt:    tc.endAt,
				Task:     mockTask,
			}

			err := job.Validate()
			if tc.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				// Check that utcStartAt and utcEndAt are correctly set
				assert.Equal(t, testTimeStart.UTC(), job.utcStartAt)
				assert.Equal(t, testTimeEnd.UTC(), job.utcEndAt)

				// Re-run Validate to ensure utcStartAt and utcEndAt do not change
				initialStartAt := job.utcStartAt
				initialEndAt := job.utcEndAt
				err = job.Validate()
				assert.NoError(t, err)
				assert.Equal(t, initialStartAt, job.utcStartAt)
				assert.Equal(t, initialEndAt, job.utcEndAt)
			}

			// 2nd pass to ensure Validate doesn't increment times x2
			err = job.Validate()
			if tc.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				// Check that utcStartAt and utcEndAt are correctly set
				assert.Equal(t, testTimeStart.UTC(), job.utcStartAt)
				assert.Equal(t, testTimeEnd.UTC(), job.utcEndAt)

				// Re-run Validate to ensure utcStartAt and utcEndAt do not change
				initialStartAt := job.utcStartAt
				initialEndAt := job.utcEndAt
				err = job.Validate()
				assert.NoError(t, err)
				assert.Equal(t, initialStartAt, job.utcStartAt)
				assert.Equal(t, initialEndAt, job.utcEndAt)
			}
		})
	}
}

func TestJob_Validate_MultipleJSON(t *testing.T) {
	// Define an array of JSON strings with varying dates and time zones.
	jsonStrings := []string{
		`{
			"timeZone": "UTC",
			"crontab": "0 12 * * *",
			"runImmediately": true,
			"runLimit": 1,
			"startAt": null,
			"endAt": "2030-08-22T12:00:00Z",
			"task": {
				"type": "shell",
				"scriptToRun": "script2.py"
			}
		}`,
		`{
			"timeZone": "America/Chicago",
			"crontab": "0 12 * * *",
			"runImmediately": true,
			"runLimit": 1,
			"startAt": null,
			"endAt": "2030-08-22T12:00:00Z",
			"task": {
				"type": "shell",
				"scriptToRun": "script2.py"
			}
		}`,
	}

	for _, jsonString := range jsonStrings {
		var job JobPlanShell
		err := json.Unmarshal([]byte(jsonString), &job)
		assert.NoError(t, err, "Failed to unmarshal JSON")

		err = job.Validate()
		assert.NoError(t, err, "Validate failed")

		// Since StartAt is null in the sample, we only check EndAt
		expectedEndAt, _ := time.Parse(time.RFC3339, "2030-08-22T12:00:00Z")
		assert.Equal(t, expectedEndAt, job.utcEndAt, "EndAt was not set correctly")

		// Re-run Validate to ensure utcStartAt and utcEndAt do not change
		initialEndAt := job.utcEndAt
		err = job.Validate()
		assert.NoError(t, err, "Validate failed on second run")
		assert.Equal(t, initialEndAt, job.utcEndAt, "EndAt changed on second validation")
	}
}

// TestJob_SetupGoCronJob tests the SetupGoCronJob method of the Job struct.
func TestJob_SetupGoCronJob(t *testing.T) {
	jobWithCrontab := JobPlan{
		Crontab: "* * * * *",
		Task: &MockITask{
			Executed: false,
		},
	}

	_, optionsWithCrontab, err := jobWithCrontab.SetupGoCronJob()
	assert.NoError(t, err, "Expected no error for job with valid crontab")
	assert.NotEmpty(t, optionsWithCrontab, "Expected options to be set for job with crontab")

	jobWithRunLimit := JobPlan{
		RunLimit: 1,
		Task: &MockITask{
			Executed: false,
		},
	}

	_, optionsWithRunLimit, err := jobWithRunLimit.SetupGoCronJob()
	assert.NoError(t, err, "Expected no error for job with run limit")
	assert.NotEmpty(t, optionsWithRunLimit, "Expected options to be set for job with run limit")
}
