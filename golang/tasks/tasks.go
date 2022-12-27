package tasks

import (
	"errors"
	"fmt"
	"time"
)

const (
	CTFormat = time.RFC3339
	FTFormat = time.RFC3339Nano
)

// A Ttype represents a meaninglessness of our life
type Ttype struct {
	Id int
	// Время создания; можно сделать time.Time, но бывают всякие структуры, приходящие извне,
	// поэтому раз так было изначально, оставляю
	CT string
	FT string // Время выполнения; оставляю строкой для единообразия с CT
	// Саму таску и её результат можно разделить на две структуры, но опять же, зависит от контекста, поэтому оставляю
	TaskResult      string
	ProcessingError error // Ошибка выполнения корректной таски
}

// TtypeProcessor можно будет использовать для создания мока для тестирования всего обработчика,
// который сейчас находится в main (т.е. для интеграционного теста)
type TtypeProcessor interface {
	Process(*Ttype) error
}

type TtypeProcessorMain struct {
	TTL time.Duration // Время жизни таски, иначе фейл
}

var TooOldTaskErr = errors.New("the task is too old")

func (p TtypeProcessorMain) Process(task *Ttype) error {
	cT, err := time.Parse(CTFormat, task.CT)
	if err != nil {
		return NewErrCreationTime(task.CT)
	}

	if cT.After(time.Now().Add(-p.TTL)) { // ещё одно условие ошибочных тасков?
		task.TaskResult = "task has been succeeded"
		task.ProcessingError = nil
	} else {
		task.ProcessingError = TooOldTaskErr
	}
	task.FT = time.Now().Format(FTFormat)

	return nil
}

type ErrCreationTime struct {
	cT string
}

func NewErrCreationTime(cT string) ErrCreationTime {
	return ErrCreationTime{
		cT: cT,
	}
}

func (e ErrCreationTime) Error() string {
	return fmt.Sprintf("creation time incorrect format: %s", e.cT)
}

func TaskCreator(output chan<- Ttype, interval time.Duration, exit <-chan struct{}) <-chan struct{} {
	stopped := make(chan struct{})
	go func() {
		defer close(stopped)
		defer close(output)
		for {
			select {
			case <-exit:
				return
			default:
				now := time.Now()
				id := int(now.Unix())
				cT := now.Format(CTFormat)
				// Если вид ошибочных тасок существенен, то условие их появления нужно вынести в аргументы функции
				if now.Nanosecond()%2 > 0 { // условие появления ошибочных тасок: на моём компьютере всегда ложь
					cT = "Some error occurred"
				}
				output <- Ttype{CT: cT, Id: id} // передаем таск на выполнение
			}

			time.Sleep(interval)
		}
	}()

	return stopped
}
