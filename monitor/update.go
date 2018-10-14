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

// update polls the site and updates it in the database.
// The first return value indicates whether the notifier should be triggered.
func (m *Monitor) update(conn *db.Conn, s *db.Site) (bool, error) {
	var (
		err       = m.check(s)
		now       = time.Now()
		nextPoll  = now.Add(time.Duration(s.PollInterval) * time.Second)
		oldStatus = s.Status
	)
	s.LastPoll = &now
	s.NextPoll = &nextPoll
	if err == nil {
		s.Status = db.StatusUp
	} else {
		s.Status = db.StatusDown
	}
	if oldStatus != s.Status {
		s.StatusTime = &now
	}
	if err := conn.Save(s).Error; err != nil {
		return false, err
	}
	// Don't create / update an outage if the status was unknown
	if oldStatus == s.Status || oldStatus == db.StatusUnknown {
		return false, nil
	}
	o := &db.Outage{}
	switch s.Status {
	case db.StatusDown:
		m.log.Infof("%s is offline", s.Name)
		o.StartTime = now
		o.Description = err.Error()
		o.SiteID = s.ID
	case db.StatusUp:
		m.log.Infof("%s is back online", s.Name)
		if db := conn.
			Model(o).
			Where("end_time IS NULL AND site_id = ?", s.ID).
			Update("end_time", &now); db.Error != nil {
			if db.RecordNotFound() {
				return false, nil
			}
		}
	}
	return true, conn.Save(o).Error
}
