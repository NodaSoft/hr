package main

import (
	"time"
)

// ЗАДАНИЕ:
// * сделать из плохого кода хороший;
// * важно сохранить логику появления ошибочных тасков;
// * сделать правильную мультипоточность обработки заданий.
// Обновленный код отправить через merge-request.

// приложение эмулирует получение и обработку тасков, пытается и получать и обрабатывать в многопоточном режиме
// В конце должно выводить успешные таски и ошибки выполнены остальных тасков

// TODO: добавить контекст и обработку сигналов
func main() {

	// TODO: возможно стоит увеличить буффер, иначе taskSorter будет периодически застревать на записи в канал
	doneTasks := make(chan Ttype)
	undoneTasks := make(chan error)

	// go func() {
	// 	// получение тасков
	// 	for t := range superChan {
	// 		// синхронно получаем таск, асинхронно обрабатываем. Может наоборот?
	// 		t = taskWorker(t)
	// 		go taskSorter(t, doneTasks, undoneTasks)
	// 	}
	// 	// TODO: мы нигде не закрываем superChan, поэтому цикл будет работать бесконечно и мы никогда не попадем на эту строчку
	// 	close(superChan) // косвенно описывает режим работы программы: дожидаемся выполнения всех тасков и только потом разбираем результаты
	// }()

	result := map[int]models.Ttype{}
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
