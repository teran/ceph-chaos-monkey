package monkey

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestStats(t *testing.T) {
	r := require.New(t)

	stats := NewStats()

	// Test ObserveWrite
	stats.ObserveWrite(100*time.Millisecond, nil)
	stats.ObserveWrite(200*time.Millisecond, errors.New("write error"))

	// Test ObserveRead
	stats.ObserveRead(150*time.Millisecond, nil)
	stats.ObserveRead(250*time.Millisecond, errors.New("read error"))

	// Dump stats and validate
	result := stats.Dump()

	r.Equal(MeasurementValue{
		AvgWritesLatency:     150 * time.Millisecond,
		AvgReadsLatency:      200 * time.Millisecond,
		WritesCountTotal:     2,
		WritesErrorsTotal:    1,
		WritesSuccessPercent: 0.50,
		ReadsCountTotal:      2,
		ReadsErrorsTotal:     1,
		ReadsSuccessPercent:  0.50,
	}, result)
}
