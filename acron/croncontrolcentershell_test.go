package acron

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"
)

type JobPlanShell struct {
	JobPlan
	CDRootDir string `json:"cdRootDir,omitempty"`

	// Turn off logging
	TurnOffLogs bool `json:"turnOffLogs,omitempty"`
	// If IsLogsInsideScriptDir is true, then output logs
	IsLogsInsideScriptDir bool `json:"isLogsInsideScriptDir,omitempty"`

	TestBool bool `json:"testBool,omitempty"`

	mu sync.RWMutex
}

func (j *JobPlanShell) GetCDRootDir() string {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return j.CDRootDir
}

func (j *JobPlanShell) SetFilePath(filePath string) {
	j.mu.Lock()
	defer j.mu.Unlock()
	if j.CDRootDir == "" {
		j.CDRootDir = path.Dir(filePath)
	}
	j.filePath = filePath
}

func (j *JobPlanShell) GetRunFunction() (function any) {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return runAppJobPlan
}

// UnmarshalJSON is a custom unmarshaller for JobPlan that handles ITask.
func (j *JobPlanShell) UnmarshalJSON(data []byte) error {
	j.mu.Lock()
	defer j.mu.Unlock()

	type Alias JobPlanShell
	aux := &struct {
		Task json.RawMessage `json:"task"`
		*Alias
	}{
		Alias: (*Alias)(j),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return fmt.Errorf("failed to unmarshal JobPlan: %v", err)
	}

	if err := j.UnmarshalJSONTask(aux.Task); err != nil {
		return fmt.Errorf("failed to unmarshal JobPlanTask: %v", err)
	}

	return nil
}

// runAppJobPlan runs the provided IJobPlan interface implementation.
// Required by gocron.
func runAppJobPlan(jobPlan IJobPlan) {
	_, _ = jobPlan.Run(nil)
}

var muDefaultLogFileTest sync.Mutex
var testDetectLogFiles int
var testDetectLogFileErrs int

func (j *JobPlanShell) Run(_ ICronControlCenter) (IJRun, error) {
	// Ignore ccc passed into function
	ccc := &CronControlCenterShell{
		CronControlCenter: CronControlCenter{},
		HideStdOut:        true,
		CDRootDir:         j.GetCDRootDir(),
		LogStds:           LogStds{},
	}
	ccc.SetJRun(NewJRunWithOptions(j.GetJobPlanId(), j.GetTitle(), j.GetTask().GetType()))

	jrun, err := j.RunJobPlanDefault(ccc)

	if !j.TurnOffLogs {
		if err := ccc.SaveLogs(false, path.Join(j.GetCDRootDir(), "logs")); err != nil {
			return jrun, fmt.Errorf("failed to save logs: %v", err)
		}
		// Save applies the index as the filename but not apply a timestamp.
		// In production, change this to use SaveWithTimestamp or other filename.
		muDefaultLogFileTest.Lock()
		testDetectLogFiles++
		if ccc.GetJRun().GetError() != nil {
			testDetectLogFileErrs++
		} else {
			for _, logStd := range ccc.GetLogStds() {
				if !strings.HasPrefix(logStd.StdOut, "Hello World from script") {
					testDetectLogFileErrs++
				}
			}
		}
		muDefaultLogFileTest.Unlock()
	}

	return jrun, err
}

// TestStartScheduler_JobPlanShell tests starting the scheduler and its state.
func TestStartScheduler_JobPlanShell(t *testing.T) {
	if err := SetScheduler(nil, true); err != nil {
		t.Fatal(err)
	}

	// Assuming the test_data directory is in the same directory as the test file.
	workingDir := "test_data"
	if err := removeLogsDirs(workingDir); err != nil {
		t.Error(err)
		return
	}

	// Call LoadJobJSONFiles with the reflect.Type of JobPlanShell
	jobs, err := LoadJobJSONFiles(workingDir, reflect.TypeOf(JobPlanShell{}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify the loaded jobs
	if len(jobs) != 2 {
		t.Fatalf("expected 2 jobs, got %d", len(jobs))
	}

	testDetectLogFiles = 0
	testDetectLogFileErrs = 0

	for ii, job := range jobs {
		j, ok := job.(*JobPlanShell)
		if !ok {
			t.Errorf("expected JobPlanShell, got %T at %d", j, ii)
			return
		}
		assert.Equal(t, j.TestBool, true)
		assert.True(t, strings.Contains(j.GetCDRootDir(), "test_data/script"))
		err := AddJobPlan(job)
		assert.NoError(t, err, "Expected no error when adding a job")
		assert.Equal(t, ii+1, len(SCHEDULER().Jobs()))
	}

	SCHEDULER().Start()
	defer func() {
		err := SCHEDULER().StopJobs()
		assert.NoError(t, err, "Expected no error when stopping jobs")
	}()

	assert.Equal(t, len(jobs), len(SCHEDULER().Jobs()))

	// Polling loop with timeout to check the global job execution status.
	timeout := time.After(3 * time.Second)
	tick := time.Tick(200 * time.Millisecond)

	for {
		select {
		case <-timeout:
			t.Fatal("Test timed out: job execution status was not set to true within 3 seconds")
		case <-tick:
			muDefaultLogFileTest.Lock()
			myCount := testDetectLogFiles
			myCountErrs := testDetectLogFileErrs
			muDefaultLogFileTest.Unlock()
			if myCount == 2 {
				assert.Equal(t, 0, myCountErrs)
				return
			}
		}
	}
}

// removeLogsDirs finds and removes any directory labeled "logs" under the specified root directory.
func removeLogsDirs(root string) error {
	// Use Glob to find all "logs" directories under the root directory
	pattern := filepath.Join(root, "*", "logs")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil // fmt.Errorf("failed to find logs directories: %v", err)
	}

	// Iterate through the matches and remove each "logs" directory
	for _, match := range matches {
		//fmt.Printf("Removing directory: %s\n", match)
		if err := os.RemoveAll(match); err != nil {
			return fmt.Errorf("failed to remove directory %s: %v", match, err)
		}
	}

	return nil
}
