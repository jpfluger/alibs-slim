package aconns

import (
	"database/sql" // Package sql provides a generic interface around SQL (or SQL-like) databases.
	"fmt"
	"time"
)

// ISBAdapterSqlRows defines the interface for SQL rows with timing functionality.
type ISBAdapterSqlRows interface {
	StartTimer()
	EndTimer() time.Duration
	GetTimerDuration() time.Duration
	Next() bool
	Close() error
	Scan(dest ...any) error
}

// SBAdapterSqlRows implements ISBAdapterSqlRows for sandbox environments.
type SBAdapterSqlRows struct {
	rows      *sql.Rows
	timeStart time.Time
	timeEnd   time.Time
}

// StartTimer starts the timer for the SQL rows operation.
func (rows *SBAdapterSqlRows) StartTimer() {
	rows.timeStart = time.Now()
	rows.timeEnd = time.Time{}
}

// EndTimer ends the timer for the SQL rows operation and returns the duration.
func (rows *SBAdapterSqlRows) EndTimer() time.Duration {
	rows.timeEnd = time.Now()
	return rows.GetTimerDuration()
}

// GetTimerDuration returns the duration between the start and end times.
func (rows *SBAdapterSqlRows) GetTimerDuration() time.Duration {
	if rows.timeStart.IsZero() || rows.timeEnd.IsZero() {
		return 0
	}
	return rows.timeEnd.Sub(rows.timeStart)
}

// Next advances to the next row in the result set.
func (rows *SBAdapterSqlRows) Next() bool {
	defer func() {
		if r := recover(); r != nil {
			rows.Close()
		}
	}()
	if rows.rows == nil {
		return false
	}
	if rows.timeStart.IsZero() {
		rows.StartTimer()
	}
	doNext := rows.rows.Next()
	if !doNext {
		_ = rows.rows.Close()
		rows.EndTimer()
	}
	return doNext
}

// Close closes the SQL rows.
func (rows *SBAdapterSqlRows) Close() (err error) {
	defer func() {
		if r := recover(); r != nil {
			if err == nil {
				err = fmt.Errorf("panic occurred: %v", r)
			}
		}
	}()
	if rows.rows == nil {
		return nil
	}
	return rows.rows.Close()
}

// Scan copies the columns in the current row into the values pointed at by dest.
func (rows *SBAdapterSqlRows) Scan(dest ...any) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if err == nil {
				err = fmt.Errorf("panic occurred: %v", r)
			}
		}
	}()
	if rows.rows == nil {
		return sql.ErrNoRows
	}
	return rows.rows.Scan(dest...)
}
