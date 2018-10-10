package notifier

import (
	"time"

	"github.com/ChimeraCoder/anaconda"
	"github.com/nathan-osman/my-site-monitor/db"
	"github.com/sirupsen/logrus"
)

// Notifier sends a tweet when a site's status changes.
type Notifier struct {
	conn        *db.Conn
	api         *anaconda.TwitterApi
	log         *logrus.Entry
	triggerChan chan bool
	stopChan    chan bool
	stoppedChan chan bool
}

func (n *Notifier) run() {
	defer close(n.stoppedChan)
	defer n.log.Info("notifier stopped")
	n.log.Info("notifier started")
	for {
		var timerChan <-chan time.Time
		err := n.conn.Transaction(func(conn *db.Conn) error {
			var (
				o = &db.Outage{}
			)
			if db := conn.
				Preload("Site").
				Order("start_time").
				Where("start_notification_sent = ?", false).
				Or("end_time IS NOT NULL AND end_notification_sent = ?", false).
				First(o); db.Error != nil {
				if !db.RecordNotFound() {
					return db.Error
				}
				return nil
			}
			return n.tweet(conn, o)
		})
		if err != nil {
			n.log.Errorf("%s - retrying in 30s", err)
			timerChan = time.After(30 * time.Second)
		}
		select {
		case <-timerChan:
		case <-n.triggerChan:
		case <-n.stopChan:
			return
		}
	}
}

// New creates a new notifier.
func New(cfg *Config) *Notifier {
	anaconda.SetConsumerKey(cfg.ConsumerKey)
	anaconda.SetConsumerSecret(cfg.ConsumerSecret)
	n := &Notifier{
		conn:        cfg.Conn,
		api:         anaconda.NewTwitterApi(cfg.AccessToken, cfg.AccessSecret),
		log:         logrus.WithField("context", "notifier"),
		triggerChan: make(chan bool, 1),
		stopChan:    make(chan bool),
		stoppedChan: make(chan bool),
	}
	go n.run()
	return n
}

// Trigger indicates that a site has changed status.
func (n *Notifier) Trigger() {
	select {
	case n.triggerChan <- true:
	default:
	}
}

// Close shuts down the notifier.
func (n *Notifier) Close() {
	close(n.stopChan)
	<-n.stoppedChan
}
