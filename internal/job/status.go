package job

import (
	"sync"
)

type Status string

const (
	StatusQueued  Status = "queued"
	StatusRunning Status = "running"
	StatusDone    Status = "done"
	StatusFailed  Status = "failed"
	StatusTimedOut Status = "timed_out"
	StatusUnknown Status = "unknown"
)

type JobState struct {
	ID     string
	Status Status
}

var (
	jobStatusMap = make(map[string]Status)
	statusMu     sync.RWMutex
)

func SetStatus(id string, s Status) {
	statusMu.Lock()
	defer statusMu.Unlock()
	jobStatusMap[id] = s
}

func GetStatus(id string) Status {
	statusMu.RLock()
	defer statusMu.RUnlock()
	return jobStatusMap[id]
}

func GetAllStatuses() []JobState {
	statusMu.RLock()
	defer statusMu.RUnlock()
	var states []JobState
	for id, s := range jobStatusMap {
		states = append(states, JobState{ID: id, Status: s})
	}
	return states
}

