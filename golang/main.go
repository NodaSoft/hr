package main

import (
	"fmt"
	"time"
)

// Приложение эмулирует получение и обработку неких тасков. Пытается и получать, и обрабатывать в многопоточном режиме.
// Приложение должно генерировать таски 10 сек. Каждые 3 секунды должно выводить в консоль результат всех
// обработанных к этому моменту тасков (отдельно успешные и отдельно с ошибками).

// ЗАДАНИЕ: сделать из плохого кода хороший и рабочий - as best as you can.
// Важно сохранить логику появления ошибочных тасков.
// Важно оставить асинхронные генерацию и обработку тасков.
// Сделать правильную мультипоточность обработки заданий.
// Обновленный код отправить через pull-request в github
// Как видите, никаких привязок к внешним сервисам нет - полный карт-бланш на модификацию кода.

type TResult bool

func (tr TResult) String() string {
	if tr {
		return "task has been successed"
	}

	return "something went wrong"
}

// A Ttype represents a meaninglessness of our life
type Ttype struct {
	id     int
	cT     string // время создания
	fT     string // время выполнения
	result TResult
}

func (t *Ttype) Work() {
	tt, _ := time.Parse(time.RFC3339, t.cT)
	if tt.After(time.Now().Add(-20 * time.Second)) {
		t.result = true
	} else {
		t.result = false
	}
	t.fT = time.Now().Format(time.RFC3339Nano)

	time.Sleep(time.Millisecond * 150)
}

type TFactory chan *Ttype

func (tf *TFactory) Start() {
	timer := time.NewTimer(10 * time.Second)

	go func() {
		for {
			select {
			case <-timer.C:
				close(*tf)
				return
			default:
				ft := time.Now().Format(time.RFC3339)
				if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
					ft = "Some error occured"
				}

				*tf <- &Ttype{cT: ft, id: int(time.Now().Unix())} // передаем таск на выполнение
			}
		}
	}()
}

type TSorts struct {
	Done   chan *Ttype
	Undone chan error
}

func NewTSorts() *TSorts {
	return &TSorts{make(chan *Ttype), make(chan error)}
}

func (ts *TSorts) Sort(t *Ttype) {
	if t.result == true {
		ts.Done <- t
	} else {
		ts.Undone <- fmt.Errorf("Task id %d time %s, error %s", t.id, t.cT, t.result)
	}
}

func (ts *TSorts) Close() {
	close(ts.Done)
	close(ts.Undone)
}

type TStorage struct {
	results map[int]*Ttype
	errors  []error
}

func NewTStorage() *TStorage {
	return &TStorage{map[int]*Ttype{}, []error{}}
}

func (s *TStorage) SaveTask(t *Ttype) {
	s.results[t.id] = t
}

func (s *TStorage) SaveErr(err error) {
	s.errors = append(s.errors, err)
}

func (s *TStorage) Output() {
	fmt.Println("Errors:")
	for _, e := range s.errors {
		fmt.Println(e.Error())
	}

	fmt.Println()

	fmt.Println("Done tasks:")
	for id, t := range s.results {
		fmt.Printf("%d: %s\n", id, t.result)
	}

	fmt.Printf("\n\n")

}
func main() {
	superChan := make(TFactory, 10)
	superChan.Start()

	storage := NewTStorage()
	sorter := NewTSorts()

	done := make(chan bool)

	go func() {
		defer close(done)

		// получение тасков
		for t := range superChan {
			t.Work()
			sorter.Sort(t)
		}

		sorter.Close()
		done <- true
	}()

	go func() {
		for t := range sorter.Done {
			storage.SaveTask(t)
		}
	}()

	go func() {
		for e := range sorter.Undone {
			storage.SaveErr(e)
		}
	}()

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	output := func(i int) {
		fmt.Printf("[%d]\n", i)
		storage.Output()
	}

	for i := 1; ; i++ {
		select {
		case <-done:
			output(i)
			fmt.Println("No more tasks...")
			return
		case <-ticker.C:
			output(i)
		}
	}
}
