package monitor

import (
	"net"
	"net/http"
	"time"

	"github.com/nathan-osman/my-site-monitor/db"
	"github.com/nathan-osman/my-site-monitor/notifier"
	"github.com/rickb777/date/period"
	"github.com/sirupsen/logrus"
)

// Monitor checks the status of the sites registered in the database.
type Monitor struct {
	conn        *db.Conn
	notifier    *notifier.Notifier
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
		err := m.conn.Transaction(func(conn *db.Conn) error {
			for {
				var (
					s   = &db.Site{}
					now = time.Now()
				)
				if db := conn.
					Set("gorm:query_option", "FOR UPDATE").
					Order("next_poll NULLS FIRST").
					First(s); db.Error != nil {
					if !db.RecordNotFound() {
						return db.Error
					}
					return nil
				}
				if s.NextPoll != nil && s.NextPoll.After(now) {
					m.log.Debugf(
						"waiting %s to check %s",
						period.Between(now, *s.NextPoll).Format(),
						s.Name,
					)
					timerChan = time.After(s.NextPoll.Sub(now))
					return nil
				}
				if err := m.update(conn, s); err != nil {
					return err
				}
			}
		})
		if err != nil {
			m.log.Errorf("%s - retrying in 30s", err)
			timerChan = time.After(30 * time.Second)
		}
		select {
		case <-timerChan:
		case <-m.triggerChan:
		case <-m.stopChan:
			return
		}
	}
}

// New creates a new monitor.
func New(cfg *Config) *Monitor {
	m := &Monitor{
		conn:     cfg.Conn,
		notifier: cfg.Notifier,
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
