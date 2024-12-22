package timer

import (
	"fmt"
	"github.com/robfig/cron/v3"
)

var Timers *timerManager
var persistInfo map[string]int64

func init() {
	Timers = &timerManager{
		Job:    cron.New(),
		timers: make(map[string]*domainTimer),
	}
	persistInfo = make(map[string]int64)

	Timers.Job.Start()
}

func (t *timerManager) Print() {
	for dn, timer := range t.timers {
		fmt.Printf("%d hosts registered.", len(t.timers))
		fmt.Printf("Domain name: %s, next renewal: %s\n", dn, t.Job.Entry(timer.jobID).Next.Format("2006-01-02 15:04:05"))
	}
}
