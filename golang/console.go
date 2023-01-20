package main

//Console печатает в консоль выполненные задания
type Console struct {
	//тоже самое что и в WorkerFactory только Read only
	doneTasks   <-chan *Task
	undoneTasks <-chan error
	withErr     <-chan *Task

	//канал сообщения о завершении работы
	exitCode chan bool

	//функция обертка fmt.Println печати сообщений
	//нужна для тестов. в обычной программе собно это был бы какой-нибудь интерфейс
	//который мы подменяли бы во время тестов.
	print func(a ...any)
}

//печать сообщений из каналов
func (c *Console) Println() {
	for {
		select {
		case msg, ok := <-c.doneTasks:
			if !ok {
				return
			}
			c.print("Done:", msg)
		case msg, ok := <-c.undoneTasks:
			if !ok {
				return
			}
			c.print("Undone:", msg)
		case msg, ok := <-c.withErr:
			if !ok {
				return
			}
			c.print("Done with Err:", msg)
		case <-c.exitCode:
			return
		}
	}
}
