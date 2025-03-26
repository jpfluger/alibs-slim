package aconns

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSBAdapterSqlRows_StartTimer(t *testing.T) {
	rows := &SBAdapterSqlRows{}
	rows.StartTimer()
	assert.False(t, rows.timeStart.IsZero())
	assert.True(t, rows.timeEnd.IsZero())
}

func TestSBAdapterSqlRows_EndTimer(t *testing.T) {
	rows := &SBAdapterSqlRows{}
	rows.StartTimer()
	time.Sleep(1 * time.Millisecond)
	duration := rows.EndTimer()
	assert.False(t, rows.timeEnd.IsZero())
	assert.Greater(t, duration, time.Duration(0))
}

func TestSBAdapterSqlRows_GetTimerDuration(t *testing.T) {
	rows := &SBAdapterSqlRows{}
	assert.Equal(t, time.Duration(0), rows.GetTimerDuration())
	rows.StartTimer()
	time.Sleep(1 * time.Millisecond)
	rows.EndTimer()
	assert.Greater(t, rows.GetTimerDuration(), time.Duration(0))
}

func TestSBAdapterSqlRows_Next(t *testing.T) {
	rows := &SBAdapterSqlRows{}
	assert.False(t, rows.Next())

	mockRows := &sql.Rows{}
	rows = &SBAdapterSqlRows{rows: mockRows}
	assert.False(t, rows.Next())
}

func TestSBAdapterSqlRows_Close(t *testing.T) {
	rows := &SBAdapterSqlRows{}
	assert.NoError(t, rows.Close())

	mockRows := &sql.Rows{}
	rows = &SBAdapterSqlRows{rows: mockRows}
	assert.Error(t, rows.Close())
}

func TestSBAdapterSqlRows_Scan(t *testing.T) {
	rows := &SBAdapterSqlRows{}
	assert.Error(t, rows.Scan())

	mockRows := &sql.Rows{}
	rows = &SBAdapterSqlRows{rows: mockRows}
	assert.Error(t, rows.Scan())
}
