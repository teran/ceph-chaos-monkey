package monkey

import (
	"sync"
	"time"
)

var _ Stats = (*stats)(nil)

type MeasurementValue struct {
	AvgWritesLatency     time.Duration
	AvgReadsLatency      time.Duration
	WritesCountTotal     uint64
	WritesErrorsTotal    uint64
	WritesSuccessPercent float64
	ReadsCountTotal      uint64
	ReadsErrorsTotal     uint64
	ReadsSuccessPercent  float64
}

type Stats interface {
	Dump() MeasurementValue

	ObserveWrite(time.Duration, error)
	ObserveRead(time.Duration, error)
}

type stats struct {
	mutex *sync.RWMutex

	totalWritesLatency time.Duration
	totalReadsLatency  time.Duration

	writesCountTotal  uint64
	writesErrorsTotal uint64

	readsCountTotal  uint64
	readsErrorsTotal uint64
}

func NewStats() Stats {
	return &stats{
		mutex: &sync.RWMutex{},
	}
}

func (s *stats) Dump() MeasurementValue {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return MeasurementValue{
		AvgWritesLatency: s.totalWritesLatency / time.Duration(s.writesCountTotal),
		AvgReadsLatency:  s.totalReadsLatency / time.Duration(s.readsCountTotal),

		WritesCountTotal:     s.writesCountTotal,
		WritesErrorsTotal:    s.writesErrorsTotal,
		WritesSuccessPercent: 1.0 - (float64(s.writesErrorsTotal) / float64(s.writesCountTotal)),

		ReadsCountTotal:     s.readsCountTotal,
		ReadsErrorsTotal:    s.readsErrorsTotal,
		ReadsSuccessPercent: 1.0 - (float64(s.readsErrorsTotal) / float64(s.readsCountTotal)),
	}
}

func (s *stats) ObserveRead(latency time.Duration, err error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.totalReadsLatency += latency
	s.readsCountTotal++
	if err != nil {
		s.readsErrorsTotal++
	}
}

func (s *stats) ObserveWrite(latency time.Duration, err error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.totalWritesLatency += latency
	s.writesCountTotal++
	if err != nil {
		s.writesErrorsTotal++
	}
}
