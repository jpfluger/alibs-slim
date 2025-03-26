package acron

import (
	"fmt"
	"github.com/jpfluger/alibs-slim/autils"
	"os"
	"path"
	"strings"
	"sync"
	"time"
)

type ICronControlCenterShell interface {
	ICronControlCenter
	GetCDRootDir() string
	GetHideStdOut() bool
	SetLogStd(*LogStd)
	GetLogStds() LogStds
}

type CronControlCenterShell struct {
	CronControlCenter

	CDRootDir  string `json:"cdRootDir,omitempty"`
	HideStdOut bool   `json:"hideStdOut,omitempty"`

	LogStds `json:"logStds,omitempty"`

	mu sync.RWMutex
}

func (c *CronControlCenterShell) GetCDRootDir() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.CDRootDir
}

func (c *CronControlCenterShell) GetHideStdOut() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.HideStdOut
}

func (c *CronControlCenterShell) SetLogStd(logStd *LogStd) {
	if logStd == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.LogStds = append(c.LogStds, logStd)
}

func (c *CronControlCenterShell) GetLogStds() LogStds {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.LogStds
}

func (c *CronControlCenterShell) SaveLogs(turnOffLogs bool, dirLogs string) error {
	if turnOffLogs {
		return nil
	}

	dirLogs = strings.TrimSpace(dirLogs)
	if dirLogs == "" {
		dirLogs = c.GetCDRootDir()
		if dirLogs == "" {
			return fmt.Errorf("dirLogs is empty")
		}
		dirLogs = path.Join(dirLogs, "logs")
	}
	// In PRODUCTION, we don't want to delete the directory but do
	// ensure that it exists.
	if _, err := autils.ResolveDirectory(dirLogs); err != nil {
		if err = os.MkdirAll(dirLogs, 0755); err != nil {
			return fmt.Errorf("failed to create log directory '%s': %v", dirLogs, err)
		}
	}

	// Save applies the index as the filename but not apply a timestamp.
	// In production, change this to use SaveWithTimestamp or other filename.
	// logs/TIMESTAMP_INDEX_TYPE.txt
	// Save logs to a file
	myTime := time.Now().Format("20060102_150405")
	if err := c.GetJRun().SaveLogs(path.Join(dirLogs, fmt.Sprintf("%s_log.txt", myTime))); err != nil {
		return fmt.Errorf("failed to save jruns log: %v", err)
	}
	// Save std out/err logs to a file
	if err := c.GetLogStds().SaveStdsWithTimeStamp(dirLogs, myTime); err != nil {
		return fmt.Errorf("failed to save logs: %v", err)
	}

	return nil
}
