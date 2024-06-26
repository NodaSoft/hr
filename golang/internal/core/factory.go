package core

import (
	"time"
)

// Higher rate means exactly that , 1 makes no errors , 173 makes errors only errors
// (( actually 173 produces ~1 correct task per 362 tasks produced overall ))
// (( ofc this rate is a bit broken, but facroty is broken too , makes sense ))

var BROKEN_FACTORY_ERROR_RATE = 2

type BrokenFactory struct {
	ServiceStarted time.Time
}

func (factory *BrokenFactory) MakeTask() Task {
	now := time.Now()
	if now.Nanosecond()%BROKEN_FACTORY_ERROR_RATE > 0 {
		// Do not uncomment, will spam to the logs
		// log.Debug("Broken factory sended broken task, again, someone fix it already")
		return Task{Created: factory.ServiceStarted, Id: now.UnixNano()}
	}
	return Task{Created: now, Id: now.UnixNano()}

}

type TaskFactory interface {
	MakeTask() Task
}
