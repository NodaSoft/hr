// https://hh.ru/vacancy/73726156?hhtmFrom=vacancy_response
package main

import (
	// будет использоваться для передачи контекста в функции
	"context"
	"fmt"

	// будет использоваться для работы с мьютексами
	"sync"
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

func main() {
	// создаем контекст для отмены всех горутин
	ctx, cancel := context.WithCancel(context.Background())

	// закрываем контекст при завершении работы функции
	defer cancel()

	// Количество тасков, которые нужно обработать
	numberOfTasks := 20

	fmt.Println("Program started")
	// создаем таски (исправлена опечатка в названии)
	taskCreate := func(ctx context.Context, a chan Ttype) {
		for i := 0; i < numberOfTasks; i++ {
			// используем select для проверки контекста
			// если контекст закрыт, закрываем канал и выходим
			// иначе(т.к.) операция блокирующая, сработает default, мы запишем в канал новый таск
			select {
			case <-ctx.Done():
				fmt.Println("All tasks are done")
				close(a)
				return
			default:
				ft := time.Now().Format(time.RFC3339)
				if time.Now().Nanosecond()%2 > 0 {
					ft = "Some error occured"
				}
				fmt.Println("Task created at", ft)
				a <- Ttype{cT: ft, id: int(time.Now().Unix())}
			}
		}
		fmt.Println("All tasks are created")
		close(a) // Закрываем канал после создания заданного количества тасков
	}

	superChan := make(chan Ttype, 10)

	go taskCreate(ctx, superChan)

	task_worker := func(a Ttype) Ttype {
		tt, _ := time.Parse(time.RFC3339, a.cT)
		if tt.After(time.Now().Add(-20 * time.Second)) {
			fmt.Println("Task id", a.id, "is too fresh")
			a.taskRESULT = []byte("task has been successed")
		} else {
			fmt.Println("Task id", a.id, "is too old")
			a.taskRESULT = []byte("something went wrong")
		}
		a.fT = time.Now().Format(time.RFC3339Nano)

		time.Sleep(time.Millisecond * 150)

		return a
	}

	doneTasks := make(chan Ttype)

	// переименовал канал, чтобы было понятно, что он содержит ошибки
	erroredTasks := make(chan error)

	// Для отслеживания горутин tasksorter
	var sorterWg sync.WaitGroup
	tasksorter := func(t Ttype) {
		defer sorterWg.Done()
		if string(t.taskRESULT[14:]) == "successed" {
			doneTasks <- t
		} else {
			erroredTasks <- fmt.Errorf("Task id %d time %s, error %s", t.id, t.cT, t.taskRESULT)
		}
	}

	go func() {
		// получение тасков
		for t := range superChan {
			t = task_worker(t)
			// добавляем горутину в группу
			sorterWg.Add(1)
			go tasksorter(t)
		}
		// закрываем каналы, чтобы не было утечек
		sorterWg.Wait()
		close(doneTasks)
		close(erroredTasks)
	}()
	// создаем горутину для ожидания завершения всех горутин
	var wg sync.WaitGroup

	// создаем массив для хранения результатов
	result := make(map[int]Ttype)
	err := []error{}

	// создаем мьютекс для безопасной записи в массив
	var resultMu sync.Mutex

	// Добавляем горутину в группу
	wg.Add(1)

	// создадим две разные функции для обработки ощибок и успешных тасков
	go func() {
		// добавляем условие при завершении работы горутины
		defer wg.Done()
		for r := range doneTasks {
			// блокируем мьютекс на время записи в массив
			// если не использовать мьютекс, то возможны гонки данных
			resultMu.Lock()
			result[r.id] = r
			// разблокируем мьютекс после записи
			resultMu.Unlock()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for e := range erroredTasks {
			// блокируем мьютекс на время записи в массив
			resultMu.Lock()
			err = append(err, e)
			// разблокируем мьютекс после записи
			resultMu.Unlock()
		}
	}()

	// ожидаем завершения всех горутин
	wg.Wait()
	fmt.Println("\nProgramm reach tasks limit, trying to shutdown...\n\n")
	println("Errors:")
	for _, e := range err {
		println(e.Error())
	}

	println("Done tasks:")
	for _, r := range result {
		fmt.Printf("Task id: %d, Result: %s\n", r.id, r.taskRESULT)
	}
}
