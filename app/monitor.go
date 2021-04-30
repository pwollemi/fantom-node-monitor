package app

import (
	"time"

	"github.com/flashguru-git/node-monitor/config"
	"github.com/flashguru-git/node-monitor/log"
)

type Monitor struct {
	S    *Server
	stop chan struct{}

	Cycle   time.Duration
	RunAt   time.Time
	PrevRun time.Time

	Data map[string]interface{}
}

func NewMonitor(s *Server, c time.Duration, run time.Time, prev time.Time) *Monitor {
	return &Monitor{
		S:       s,
		stop:    make(chan struct{}),
		Cycle:   c,
		RunAt:   run,
		PrevRun: prev,
		Data: map[string]interface{}{
			"nodeId": getNodeId(),
			"ipAddr": getIpAddr(),
		},
	}
}

func (m *Monitor) Execute() {
	data := map[string]interface{}{
		"nodeId":              m.Data["nodeId"],
		"ipAddr":              m.Data["ipAddr"],
		"blockHeight":         getBlockNumber(),
		"topPeersBlockHeight": getTopPeersBlockHeight(),
		"createdAt":           time.Now(),
		"cpu":                 getCpuUsage(),
		"memory":              getMemoryUsage(),
	}

	serverURL := config.Config().GetString("SERVER_URL")
	go sendPostRequest(serverURL+"/api/nodes", data)
}

func (m *Monitor) ScheduleNextRun() {
	m.PrevRun = m.RunAt
	for m.RunAt.Before(time.Now().UTC()) {
		m.RunAt = m.RunAt.Add(m.Cycle)
	}
}

func (m *Monitor) Run() {
	for {
		var after <-chan time.Time

		remaining := m.RunAt.Sub(time.Now().UTC())
		after = time.After(remaining)

		// Sleep until the next job's run time.
		select {
		case <-after:
			log.WithFields(log.Fields{
				"DebugMessage": "sleep_finished",
			}).Debug("monitor and execute")
			m.Execute()
			m.ScheduleNextRun()
			break
		case <-m.stop:
			log.WithFields(log.Fields{
				"DebugMessage": "quit",
			}).Debug("came quit message")
			return
		}
	}
}
