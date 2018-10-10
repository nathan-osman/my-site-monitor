package monitor

import (
	"fmt"
	"time"

	"github.com/nathan-osman/my-site-monitor/db"
)

func (m *Monitor) check(s *db.Site) error {
	m.log.Debugf("checking %s...", s.Name)
	resp, err := m.client.Get(s.URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return fmt.Errorf("server returned status %s", resp.Status)
	}
	return nil
}

func (m *Monitor) update(s *db.Site) error {
	triggerNotifier := false
	err := m.conn.Transaction(func(conn *db.Conn) error {
		var (
			now       = time.Now()
			err       = m.check(s)
			oldStatus = s.Status
		)
		s.LastPoll = now
		s.NextPoll = now.Add(time.Duration(s.PollInterval) * time.Second)
		if err == nil {
			s.Status = db.StatusUp
		} else {
			s.Status = db.StatusDown
		}
		if oldStatus != s.Status {
			s.StatusTime = now
			for {
				// Don't try to update an outage if the status was unknown
				if oldStatus == db.StatusUnknown && s.Status == db.StatusUp {
					break
				}
				triggerNotifier = true
				o := &db.Outage{}
				switch s.Status {
				case db.StatusDown:
					m.log.Infof("%s is down", s.Name)
					o.StartTime = now
					o.Description = err.Error()
					o.SiteID = s.ID
				case db.StatusUp:
					m.log.Infof("%s has come back up", s.Name)
					if err := conn.
						Order("start_time DESC").
						Where("site_id = ?", s.ID).
						First(o).Error; err != nil {
						return err
					}
					o.EndTime = now
				}
				if err := conn.Save(o).Error; err != nil {
					return err
				}
				break
			}
		}
		return conn.Save(s).Error
	})
	if triggerNotifier {
		m.notifier.Trigger()
	}
	return err
}
