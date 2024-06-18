package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Вынес "магические числа" в константы
const (
	programExecutionTime = time.Second * 10
	producerTicker       = time.Millisecond * 500
	writerTicker         = time.Second * 3
	successfulTaskOffset = time.Second * -20

	taskBufferSize = 8

	writerSeparator = "\n-------------------------------------------------------\n\n"
	writerFormatStr = "Task %d, time: %s\n"
)

// Заменил Ttype на более понятное название Task
// Убрал сложные для понимания строчные переменные
type Task struct {
	id             uint64
	createdAt      time.Time
	failedToCreate bool
}

// Структура, в которой мы храним временный результат
// Воркер taskWriter каждые 3 секунды пишет содержимое в консоль и чистит элементы из слайсов
// Для избежания состояния гонки и потенциальной потери данных используем мьютекс
type TemporaryResult struct {
	mu         sync.Mutex
	successful []Task
	failed     []Task
}

func main() {
	// При достижении дедлайна контекст закрывается и рутины, его слушающие, прекращают работу
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(programExecutionTime))
	defer cancel()

	// Оставляем только один канал, в который будем помещать таски
	taskChan := make(chan Task)

	// Задаем слайсы с длиной 0 и определенной capacity.
	// Засчет этого экономим ресурсы, заранее аллоцируя память, при этом работая со слайсом как с пустым
	result := &TemporaryResult{
		successful: make([]Task, 0, taskBufferSize),
		failed:     make([]Task, 0, taskBufferSize),
	}

	// При помощи WaitGroup мы дожидаемся окончания работы всех воркеров
	wg := &sync.WaitGroup{}

	wg.Add(3)

	// Вынес функции из main для лучшей читаемости
	// producer создает таски, processor их проверяет и распределяет, а writer пишет в консоль
	go taskProducer(ctx, wg, taskChan)
	go taskProcessor(ctx, wg, taskChan, result)
	go taskWriter(ctx, wg, result)

	wg.Wait()
}

// Так как этот воркер в канал только пишет, передаем его в параметрах как write-only
// Избавился от лишнего вызова горутины, сама функция и так вызывается в рутине
// Конвертация времени в строку и обратно не служит никакой цели, избавился от нее
func taskProducer(ctx context.Context, wg *sync.WaitGroup, taskChan chan<- Task) {
	// defer хранит функции в стеке, поэтому закрытие канала помещаем ниже, чтобы канал наверняка закрылся до завершения программы
	defer wg.Done()
	defer close(taskChan)

	// Позволил себе ограничить количество создаваемых тасок, теперь таски создаются с заданной периодичностью
	ticker := time.NewTicker(producerTicker)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			timeNow := time.Now()

			// Более читаемый вид по сравнению с записью в одну строку
			task := Task{
				id:        uint64(timeNow.UnixNano()),
				createdAt: timeNow,
			}

			// Сравнение с нулем - излишняя операция, так как модуло от 2 может вернуть только 0 или 1
			// Для процессора дешевле сделать проверку на равенство
			if time.Now().Nanosecond()%2 != 0 {
				task.failedToCreate = true
			}

			taskChan <- task
		}
	}
}

// Обработчик тасков только читает из канала, поэтому канал передаем как read-only
// Структуру TemporaryResult передаем по указателю, иначе будем работать с копией и при append'ах рискуем потерять данные
func taskProcessor(ctx context.Context, wg *sync.WaitGroup, taskChan <-chan Task, result *TemporaryResult) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case task := <-taskChan:
			result.mu.Lock()

			if task.failedToCreate || task.createdAt.After(time.Now().Add(successfulTaskOffset)) {
				result.successful = append(result.successful, task)
			} else {
				result.failed = append(result.failed, task)
			}

			result.mu.Unlock()
		}
	}
}

// Каждые 3 секунды пишем в консоль накопленные таски, а потом чистим их, освобождая место под следующую итерацию
func taskWriter(ctx context.Context, wg *sync.WaitGroup, result *TemporaryResult) {
	defer wg.Done()

	ticker := time.NewTicker(writerTicker)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			result.mu.Lock()

			writeTasksToStdout(result)

			// После вывода чистим слайсы, зануляя их длину, при этом сохраняя capacity
			result.successful = result.successful[:0]
			result.failed = result.failed[:0]

			result.mu.Unlock()
		}
	}
}

// Форматирование вывода в отдельной функции, форматная строка переиспользуется, поэтому вынесена в константу
func writeTasksToStdout(result *TemporaryResult) {
	fmt.Print("Successful tasks:\n\n")

	for _, task := range result.successful {
		fmt.Printf(writerFormatStr, task.id, task.createdAt.Format(time.RFC3339))
	}

	fmt.Print("\nFailed tasks:\n\n")

	for _, task := range result.failed {
		fmt.Printf(writerFormatStr, task.id, task.createdAt.Format(time.RFC3339))
	}

	fmt.Print(writerSeparator)
}
