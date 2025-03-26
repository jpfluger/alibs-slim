package acron

import (
	"path/filepath"
	"reflect"
	"sync"
	"testing"
	"time"

	gocron "github.com/go-co-op/gocron/v2"
	"github.com/stretchr/testify/assert"
	"github.com/jpfluger/alibs-slim/areflect"
)

func init() {
	_ = areflect.TypeManager().Register(TYPEMANAGER_CRONTASKDATA, "acron-mockjobdataverify", returnTypeManagerMockJobDataVerify)
}

func returnTypeManagerMockJobDataVerify(typeName string) (reflect.Type, error) {
	var rtype reflect.Type // nil is the zero value for pointers, maps, slices, channels, and function types, interfaces, and other compound types.
	switch TaskType(typeName) {
	case TaskType("mock"):
		// Return the type of NoteFlag if typeName is "flag".
		rtype = reflect.TypeOf(MockITask{})
	}
	// Return the determined reflect.Type and no error.
	return rtype, nil
}

// MockITask is a mock implementation of the ITask interface for testing.
type MockITask struct {
	Executed bool
	mu       sync.Mutex
}

func (m *MockITask) GetType() TaskType {
	return "mock"
}

func (m *MockITask) Validate() error {
	return nil
}

func (m *MockITask) Run(ccc ICronControlCenter) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Executed = true
	return nil
}

func (m *MockITask) GetExecuted() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.Executed
}

// TestFindJobJSONFiles tests the FindJobJSONFiles function.
func TestFindJobJSONFiles(t *testing.T) {
	// Assuming the test_data directory is in the same directory as the test file.
	workingDir := "test_data"
	expectedFiles := []string{
		filepath.Join(workingDir, "script1", "job.json"),
		filepath.Join(workingDir, "script2", "job.json"),
	}

	files, err := FindJobJSONFiles(workingDir)
	assert.NoError(t, err)
	assert.ElementsMatch(t, expectedFiles, files)
}

// TestLoadJobJSONFiles tests the LoadJobJSONFiles function
func TestLoadJobJSONFiles(t *testing.T) {
	// Assuming the test_data directory is in the same directory as the test file.
	workingDir := "test_data"

	// Call LoadJobJSONFiles with the reflect.Type of MockJobPlan
	jobs, err := LoadJobJSONFiles(workingDir, reflect.TypeOf(JobPlanShell{}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify the loaded jobs
	if len(jobs) != 2 {
		t.Fatalf("expected 2 jobs, got %d", len(jobs))
	}

	if jobs[0].GetTitle() != "Job 1" || jobs[1].GetTitle() != "Job 2" {
		t.Fatalf("unexpected job titles: %v, %v", jobs[0].GetTitle(), jobs[1].GetTitle())
	}
}

// TestSCHEDULER tests the SCHEDULER function.
func TestSCHEDULER(t *testing.T) {
	scheduler := SCHEDULER()
	assert.NotNil(t, scheduler)
}

// TestSetScheduler tests the SetScheduler function.
func TestSetScheduler(t *testing.T) {
	scheduler, _ := gocron.NewScheduler(gocron.WithLocation(time.UTC))
	err := SetScheduler(scheduler, false)
	assert.NoError(t, err)
	assert.Equal(t, scheduler, globalCron)
}

// TestAddJobPlan tests the AddJobPlan function.
func TestAddJobPlan(t *testing.T) {
	jobPlan := &JobPlanShell{
		JobPlan: JobPlan{
			RunImmediately: true,
			Task:           &MockITask{},
		},
	}
	err := AddJobPlan(jobPlan)
	assert.NoError(t, err)
}

// TestScheduleJobs tests the ScheduleJobs function.
func TestScheduleJobs(t *testing.T) {
	jobPlans := IJobPlans{
		&JobPlanShell{JobPlan: JobPlan{RunImmediately: true, Task: &MockITask{}}},
		&JobPlanShell{JobPlan: JobPlan{RunImmediately: true, Task: &MockITask{}}},
	}
	err := ScheduleJobPlans(jobPlans)
	assert.NoError(t, err)
}

// TestRunITask tests the runITask function.
func TestRunITask(t *testing.T) {
	mockData := &MockITask{}
	runITask(mockData)
	// No assertion here since runITask only prints to stdout.
	// In a real-world scenario, you would use a logger that can be mocked to test the output.
}

// TestStartScheduler tests starting the scheduler and its state.
func TestStartScheduler(t *testing.T) {
	if err := SetScheduler(nil, true); err != nil {
		t.Fatal(err)
	}

	// Create a mock job.
	mockData := &MockITask{
		Executed: false,
	}
	job := &JobPlanShell{
		JobPlan: JobPlan{
			RunImmediately: true,
			// StartAt:        atime.ToPointer(time.Now().UTC().Add(1 * time.Hour)),
			Task: mockData,
		},
	}
	err := AddJobPlan(job)
	assert.NoError(t, err, "Expected no error when adding a job")
	assert.Equal(t, 1, len(SCHEDULER().Jobs()))

	SCHEDULER().Start()
	defer func() {
		err := SCHEDULER().StopJobs()
		assert.NoError(t, err, "Expected no error when stopping jobs")
	}()

	assert.Equal(t, 1, len(SCHEDULER().Jobs()))

	// Polling loop with timeout to check the global job execution status.
	timeout := time.After(3 * time.Second)
	tick := time.Tick(200 * time.Millisecond)

	for {
		select {
		case <-timeout:
			t.Fatal("Test timed out: job execution status was not set to true within 3 seconds")
		case <-tick:
			if mockData.GetExecuted() {
				assert.True(t, mockData.GetExecuted(), "The job was executed")
				return
			}
		}
	}
}
