package main

import (
	"fmt"
	"time"
)

// ЗАДАНИЕ:
// * сделать из плохого кода хороший;
// * важно сохранить логику появления ошибочных тасков;
// * сделать правильную мультипоточность обработки заданий.
// Обновленный код отправить через merge-request.

// приложение эмулирует получение и обработку тасков, пытается и получать и обрабатывать в многопоточном режиме
// В конце должно выводить успешные таски и ошибки выполнены остальных тасков

// A Ttype represents a meaninglessness of our life
type Ttype struct {
	id         int
	cT         string // время создания
	fT         string // время выполнения
	taskRESULT []byte
}

func taskCreturer(a chan Ttype) {
	go func() {
		for {
			// TODO: variables name ft -> ct
			ct := time.Now().Format(time.RFC3339)
			if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
				ct = "Some error occured"
			}
			a <- Ttype{cT: ct, id: int(time.Now().Unix())} // передаем таск на выполнение
		}
	}()
}

// emulate complete time
func taskWorker(a Ttype) Ttype {
	// TODO: handle error
	tt, _ := time.Parse(time.RFC3339, a.cT)
	// TODO: на практике всегда будет true, поскольку taskCreturer будет производить таски мгновенно
	// но сказано логику возникновения ошибки не менять
	if tt.After(time.Now().Add(-20 * time.Second)) {
		a.taskRESULT = []byte("task has been successed")
	} else {
		a.taskRESULT = []byte("something went wrong")
	}
	a.fT = time.Now().Format(time.RFC3339Nano)

	time.Sleep(time.Millisecond * 150)

	return a
}

// puts tasks into boxes depends on result.
func taskSorter(t Ttype, doneTasks chan Ttype, undoneTasks chan error) {
	// TODO: потенциальная паника, необходимо проверять длину слайса перед указанием индекса, либо поменять на bytes.Contains
	if string(t.taskRESULT[14:]) == "successed" {
		doneTasks <- t
	} else {
		undoneTasks <- fmt.Errorf("Task id %d time %s, error %s", t.id, t.cT, t.taskRESULT)
	}
}

// TODO: добавить контекст и обработку сигналов
func main() {

	superChan := make(chan Ttype, 10)

	go taskCreturer(superChan)
	// TODO: возможно стоит увеличить буффер, иначе taskSorter будет периодически застревать на записи в канал
	doneTasks := make(chan Ttype)
	undoneTasks := make(chan error)

	go func() {
		// получение тасков
		for t := range superChan {
			// синхронно получаем таск, асинхронно обрабатываем. Может наоборот?
			t = taskWorker(t)
			go taskSorter(t, doneTasks, undoneTasks)
		}
		// TODO: мы нигде не закрываем superChan, поэтому цикл будет работать бесконечно и мы никогда не попадем на эту строчку
		close(superChan) // косвенно описывает режим работы программы: дожидаемся выполнения всех тасков и только потом разбираем результаты
	}()

	result := map[int]Ttype{}
	err := []error{}
	go func() {
		// с чего мы взяли что все таски завершены
		// TODO: неплохо бы дождаться завершения всей очереди тасков
		for r := range doneTasks {
			r := r
			go func() {
				// TODO: mutex
				result[r.id] = r
			}()
		}
		// TODO: same
		for r := range undoneTasks {
			r := r
			go func() {
				// TODO: mutex
				err = append(err, r)
			}()
		}

		close(doneTasks)
		close(undoneTasks)
	}()
	// TODO: заменить на waitgroup
	time.Sleep(time.Second * 3)

	println("Errors:")
	for r := range err {
		// TODO: публиковать содержание ошибки, а не индекс в слайсе
		println(r)
	}

	println("Done tasks:")
	for r := range result {
		// TODO: same
		println(r)
	}
}
