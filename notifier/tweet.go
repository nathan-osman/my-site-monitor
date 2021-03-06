package notifier

import (
	"fmt"

	"github.com/nathan-osman/my-site-monitor/db"
	"github.com/rickb777/date/period"
)

func (n *Notifier) tweet(conn *db.Conn, o *db.Outage) error {
	var (
		status    string
		suffix    = fmt.Sprintf("\n\nOutage ID: %d", o.ID)
		suffixLen = len(suffix)
	)
	if !o.StartNotificationSent {
		status = fmt.Sprintf(
			"%s is offline - %s",
			o.Site.Name,
			o.Description,
		)
		o.StartNotificationSent = true
	} else {
		status = fmt.Sprintf(
			"%s is back online - site was offline for %s",
			o.Site.Name,
			period.Between(o.StartTime, *o.EndTime).Format(),
		)
		o.EndNotificationSent = true
	}
	if err := conn.Save(o).Error; err != nil {
		return err
	}
	if len(status) > 280-suffixLen {
		status = status[:suffixLen-1] + "…"
	}
	status += suffix
	n.log.Infof("sending tweet for %s", o.Site.Name)
	_, err := n.api.PostTweet(status, nil)
	return err
}
