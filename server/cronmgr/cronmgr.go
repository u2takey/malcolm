package cronmgr

import (
	"sort"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"

	"github.com/robfig/cron"

	"github.com/arlert/malcolm/model"
)

// Interface ...
type Interface interface {
	Start() error
	UpInsert(pipe *model.Pipeline) error
	Delete(pipeid string) error
	CronPipeChan() chan string
}

// Schedule ...
type Schedule interface {
	Next(time.Time) time.Time
}

// Entry ...
type Entry struct {
	Schedule    Schedule
	ScheduleStr string
	Next        time.Time
	Prev        time.Time
	Pipeid      string
}

// CronMgr ..
type CronMgr struct {
	entrylock sync.Mutex
	entries   []*Entry // pipeid : schedule
	pipechan  chan string
	update    chan *Entry
	running   bool
}

// New ...
func New() *CronMgr {
	return &CronMgr{
		pipechan: make(chan string),
		update:   make(chan *Entry),
	}
}

// CronPipeChan ...
func (m *CronMgr) CronPipeChan() chan string {
	return m.pipechan
}

// Start running
func (m *CronMgr) Start() error {
	m.running = true
	go m.run()
	return nil
}

// Delete entry
func (m *CronMgr) Delete(pipeid string) error {
	entry := &Entry{
		Pipeid: pipeid,
	}
	m.update <- entry
	return nil
}

// UpInsert canbe update/insert/delete entry to entries
func (m *CronMgr) UpInsert(pipe *model.Pipeline) error {
	entry, err := m.entryWithPipe(pipe)
	if err != nil {
		return err
	}
	m.update <- entry
	return nil
}

func (m *CronMgr) entryWithPipe(pipe *model.Pipeline) (entry *Entry, err error) {
	entry = &Entry{
		Pipeid: pipe.ID.Hex(),
	}
	if pipe.Trigger.Cron != nil {
		entry.ScheduleStr = pipe.Trigger.Cron.Schedule
		sche, err1 := cron.Parse(entry.ScheduleStr)
		if err1 != nil {
			err = err1
			return
		}
		entry.Schedule = sche
	}
	return
}

func (m *CronMgr) run() {
	now := time.Now()
	for _, entry := range m.entries {
		entry.Next = entry.Schedule.Next(now)
	}

	for {
		sort.Sort(byTime(m.entries))

		var timer *time.Timer
		if len(m.entries) == 0 || m.entries[0].Next.IsZero() {
			timer = time.NewTimer(100000 * time.Hour)
		} else {
			timer = time.NewTimer(m.entries[0].Next.Sub(now))
		}

		for {
			select {
			case now = <-timer.C:
				now = time.Now()
				for _, e := range m.entries {
					if e.Next.After(now) || e.Next.IsZero() {
						break
					}
					m.pipechan <- e.Pipeid
					e.Prev = e.Next
					e.Next = e.Schedule.Next(now)
				}

			case newentry := <-m.update:
				timer.Stop()
				now = time.Now()
				action := -999
				for index, entry := range m.entries {
					if entry.Pipeid == newentry.Pipeid {
						if newentry.Schedule == nil {
							//delete
							action = index
						} else {
							//update
							action = -1
							entry.Schedule = newentry.Schedule
							entry.Next = entry.Schedule.Next(now)
						}
						break
					}
				}
				if action >= 0 {
					// delete
					m.entries = append(m.entries[:action], m.entries[action+1:]...)
				}
				if action < -1 && newentry.Schedule != nil {
					// insert
					newentry.Next = newentry.Schedule.Next(now)
					m.entries = append(m.entries, newentry)
				}
				logrus.Debug("pipe in cron mgr after update ", len(m.entries))
			}
			break
		}
	}
}

type byTime []*Entry

func (s byTime) Len() int      { return len(s) }
func (s byTime) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s byTime) Less(i, j int) bool {
	if s[i].Next.IsZero() {
		return false
	}
	if s[j].Next.IsZero() {
		return true
	}
	return s[i].Next.Before(s[j].Next)
}
