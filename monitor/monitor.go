package monitor

import (
	"net"
	"net/http"
	"time"

	"github.com/nathan-osman/my-site-monitor/db"
	"github.com/sirupsen/logrus"
)

// Monitor checks the status of the sites registered in the database.
type Monitor struct {
	conn        *db.Conn
	client      *http.Client
	log         *logrus.Entry
	triggerChan chan bool
	stopChan    chan bool
	stoppedChan chan bool
}

func (m *Monitor) run() {
	defer close(m.stoppedChan)
	defer m.log.Info("monitor stopped")
	m.log.Info("monitor started")
	for {
		var timerChan <-chan time.Time
		for {
			var (
				s   = &db.Site{}
				now = time.Now()
			)
			if db := m.conn.Order("next_poll").First(s); db.Error != nil {
				if !db.RecordNotFound() {
					m.log.Errorf("%s - retrying in 30s", db.Error.Error())
					timerChan = time.After(30 * time.Second)
				}
				break
			}
			if s.NextPoll.After(now) {
				d := s.NextPoll.Sub(now)
				m.log.Debugf("waiting %s to check %s", d, s.Name)
				timerChan = time.After(d)
				break
			}
			if err := m.update(s); err != nil {
				m.log.Errorf("%s", err.Error())
			}
		}
		select {
		case <-timerChan:
		case <-m.stopChan:
			return
		}
	}
}

// New creates a new monitor.
func New(cfg *Config) *Monitor {
	m := &Monitor{
		conn: cfg.Conn,
		client: &http.Client{
			Transport: &http.Transport{
				Dial: (&net.Dialer{
					Timeout: 5 * time.Second,
				}).Dial,
			},
			Timeout: 5 * time.Second,
		},
		log:         logrus.WithField("context", "monitor"),
		triggerChan: make(chan bool, 1),
		stopChan:    make(chan bool),
		stoppedChan: make(chan bool),
	}
	go m.run()
	return m
}

// Trigger indicates that the database has been modified.
func (m *Monitor) Trigger() {
	select {
	case m.triggerChan <- true:
	default:
	}
}

// Close shuts down the monitor.
func (m *Monitor) Close() {
	close(m.stopChan)
	<-m.stoppedChan
}
