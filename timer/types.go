package timer

import (
	"auto_cert/manager"
	"github.com/robfig/cron/v3"
)

type domainTimer struct {
	domain *manager.Domain
	jobID  cron.EntryID
}
type timerManager struct {
	Job    *cron.Cron
	timers map[string]*domainTimer
}
