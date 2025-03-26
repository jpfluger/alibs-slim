package acron

import (
	"errors"
	"fmt"
	"github.com/gofrs/uuid/v5"
	"github.com/rs/zerolog"
	"github.com/jpfluger/alibs-slim/aerr"
	"os"
	"sync"
	"time"
)

// IJRun interface defines the methods that JRun must implement.
type IJRun interface {
	GetJobPlanId() uuid.UUID
	GetJobPlanTitle() string
	GetTaskType() TaskType
	Begin()
	End()
	IsFinished() bool
	GetError() error
	SetError(error)
	GetStartTime() time.Time
	GetEndTime() *time.Time
	GetLogs() []string
	SaveLogs(filePath string) error
	Logger() *zerolog.Logger
}

type IJRuns []IJRun

// JRun struct holds the details of a cron job run.
type JRun struct {
	Error        *aerr.Error    `json:"error,omitempty"`   // Error encountered during the job run
	StartTime    time.Time      `json:"startTime"`         // Time when the cron job started
	EndTime      *time.Time     `json:"endTime,omitempty"` // Time when the cron job ended (pointer to allow nil value)
	Logs         []string       `json:"logs,omitempty"`    // Slice to store log events
	logger       zerolog.Logger // Logger instance
	JobPlanId    uuid.UUID      `json:"jobPlanId"`    // Job plan id at the time of the run
	JobPlanTitle string         `json:"jobPlanTitle"` // Job plan title at the time of the run
	TaskType     TaskType       `json:"taskType"`     // TaskType at the time of the run
	mu           sync.RWMutex   // Mutex for concurrency safety
	logMu        sync.Mutex     // Separate mutex for logging
}

// logJRunWriter is an io.Writer that writes to the Logs slice.
type logJRunWriter struct {
	j *JRun
}

// Write implements the io.Writer interface for logJRunWriter.
// It appends the log message to the Logs slice of JRun.
func (lw *logJRunWriter) Write(p []byte) (n int, err error) {
	lw.j.logMu.Lock()
	defer lw.j.logMu.Unlock()
	lw.j.Logs = append(lw.j.Logs, string(p))
	return len(p), nil
}

// NewJRun creates a new JRun instance with an initialized logger.
// It returns a pointer to the newly created JRun.
func NewJRun() *JRun {
	return NewJRunWithOptions(uuid.Nil, "", "")
}

// NewJRunWithOptions creates a new JRun instance with specified job plan ID, title and taskType.
// It returns a pointer to the newly created JRun.
func NewJRunWithOptions(jobPlanId uuid.UUID, jobPlanTitle string, taskType TaskType) *JRun {
	j := &JRun{
		Logs:         []string{},
		JobPlanId:    jobPlanId,
		JobPlanTitle: jobPlanTitle,
		TaskType:     taskType,
	}
	j.logger = zerolog.New(&logJRunWriter{j: j}).With().Timestamp().Str("jobPlanId", jobPlanId.String()).Logger()
	return j
}

// Logger returns the zerolog.Logger associated with the JRun instance.
func (j *JRun) Logger() *zerolog.Logger {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return &j.logger
}

// Begin marks the start time of the cron job.
func (j *JRun) Begin() {
	j.mu.Lock()
	defer j.mu.Unlock()
	j.StartTime = time.Now()
	j.logger.Info().Msg("job started")
}

var ErrStartTimeHasNoValue = errors.New("start time has no value")

// End marks the end time of the cron job.
// When start time has no value, an error is created and
// StartTime gets the EndTime value.
func (j *JRun) End() {
	j.mu.Lock()
	defer j.mu.Unlock()
	now := time.Now().UTC() // Get current UTC time
	if j.StartTime.IsZero() {
		if j.Error == nil {
			j.Error = aerr.NewError(ErrStartTimeHasNoValue)
		} else {
			j.Error = aerr.NewError(fmt.Errorf("%v; %v", ErrStartTimeHasNoValue, j.Error.Error()))
		}
		j.StartTime = now
	}
	j.EndTime = &now // Set EndTime to the pointer of the current time
	j.logger.Info().Msg("job ended")
}

// GetJobPlanId returns the job plan ID associated with the cron job.
func (j *JRun) GetJobPlanId() uuid.UUID {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return j.JobPlanId
}

// GetJobPlanTitle returns the job plan title associated with the cron job.
func (j *JRun) GetJobPlanTitle() string {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return j.JobPlanTitle
}

// GetTaskType returns the job task type associated with the cron job.
func (j *JRun) GetTaskType() TaskType {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return j.TaskType
}

// IsFinished returns true if the cron job has finished.
func (j *JRun) IsFinished() bool {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return j.EndTime != nil
}

// GetError returns the error encountered during the job run, if any.
func (j *JRun) GetError() error {
	j.mu.RLock()
	defer j.mu.RUnlock()
	if j.Error != nil {
		return j.Error.ToError()
	}
	return nil
}

// SetError sets the error encountered during the job run.
func (j *JRun) SetError(err error) {
	j.mu.Lock()
	defer j.mu.Unlock()
	j.Error = aerr.NewError(err)
}

// GetStartTime returns the start time of the cron job.
func (j *JRun) GetStartTime() time.Time {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return j.StartTime
}

// GetEndTime returns the end time of the cron job.
func (j *JRun) GetEndTime() *time.Time {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return j.EndTime
}

// GetLogs returns the logs captured during the cron job run.
func (j *JRun) GetLogs() []string {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return j.Logs
}

// SaveLogs saves the logs of the JRun to the specified file path.
func (j *JRun) SaveLogs(filePath string) error {
	// Create or open the log file for writing.
	logFile, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create log file %s: %v", filePath, err)
	}
	defer logFile.Close()

	// Write each log entry to the file.
	for _, log := range j.GetLogs() {
		if _, err := logFile.WriteString(log); err != nil {
			return fmt.Errorf("failed to write log to file %s: %v", filePath, err)
		}
	}
	return nil
}

// GetByJobPlanId returns a sublist of IJRuns with the specified JobPlanId.
func (runs IJRuns) GetByJobPlanId(jobPlanId uuid.UUID) IJRuns {
	var result IJRuns
	for _, run := range runs {
		if run.GetJobPlanId() == jobPlanId {
			result = append(result, run)
		}
	}
	return result
}

// GetByJobPlanTitle returns a sublist of IJRuns with the specified JobPlanTitle.
func (runs IJRuns) GetByJobPlanTitle(jobPlanTitle string) IJRuns {
	var result IJRuns
	for _, run := range runs {
		if run.GetJobPlanTitle() == jobPlanTitle {
			result = append(result, run)
		}
	}
	return result
}

// GetByTaskType returns a sublist of IJRuns with the specified TaskType.
func (runs IJRuns) GetByTaskType(taskType TaskType) IJRuns {
	var result IJRuns
	for _, run := range runs {
		if run.GetTaskType() == taskType {
			result = append(result, run)
		}
	}
	return result
}

// GetFinished returns a sublist of IJRuns that have finished.
func (runs IJRuns) GetFinished() IJRuns {
	var result IJRuns
	for _, run := range runs {
		if run.IsFinished() {
			result = append(result, run)
		}
	}
	return result
}

// GetNotFinished returns a sublist of IJRuns that have not finished.
func (runs IJRuns) GetNotFinished() IJRuns {
	var result IJRuns
	for _, run := range runs {
		if !run.IsFinished() {
			result = append(result, run)
		}
	}
	return result
}

// ReplaceByJobPlanId replaces an existing IJRun with the same JobPlanId.
func (runs *IJRuns) ReplaceByJobPlanId(newRun IJRun) {
	for i, run := range *runs {
		if run.GetJobPlanId() == newRun.GetJobPlanId() {
			(*runs)[i] = newRun
			return
		}
	}
	*runs = append(*runs, newRun)
}
