package main

import (
	"sync"
	"time"
)

// Postman создает задачи. ака эмулирует сообщения из RabbitMQ
// это измененная изначальная функция create, теперь с условием завершения программы
// ибо изначально условия завершения программы были нечетко определены.
type Postman struct {
	//mu не дает создать новый тикер пока работает старый
	mu sync.Mutex
	//ticker для генерации сообщений
	ticker *time.Ticker
	//задержка до генерации
	delay    time.Duration
	//время работы
	schedule time.Duration
	//stop канал для остановки генерации сообщений
	stop chan bool
	//счетчик созданых сообщений
	numberCreated int
}

func (c *Postman) Create() <-chan *Task {
	c.mu.Lock()
	c.ticker = time.NewTicker(c.delay)

	go func() {
		time.Sleep(c.schedule)
		c.stop <- true
	}()

	//Изначально было 10, но 100 мне нравится больше.
	//Иначе таски не попадают в undone т.к канал всегда полный
	tasks := make(chan *Task, 100)
	go c.do(tasks)
	return tasks
}

func (c *Postman) do(tasks chan *Task) {
	//где-то в go 1.16 уже поправили накладные работы на скорость работы defer и сильно снизили их со 100мс.
	//так что сейчас ок использовать defer и mutex.Unlock
	defer c.mu.Unlock()
	for {
		select {
		case <-c.stop:
			c.ticker.Stop()
			close(tasks)
			return
		case <-c.ticker.C:
			tasks <- NewTask(time.Now())
			c.numberCreated++
		}
	}
}

// NewPostman создает новый Postman
// schedule время работы до остановки
// delay время между генерацией сообщений
func NewPostman(schedule, delay time.Duration) *Postman {
	stop := make(chan bool, 1)

	return &Postman{
		delay:    delay,
		schedule: schedule,
		stop:     stop,
	}
}
