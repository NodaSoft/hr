package main

import (
	"github.com/sirupsen/logrus"
	"io"
	"taskService/config"
	"taskService/taskEmitter"
	"taskService/taskRouter"
	"taskService/taskWorker"
	"time"
)

func main() {
	// загружаем аргументы запуска
	// TODO: можно добавить туда загрузку конфига в yaml\json
	config.Init()

	// настраиваем глобальный логгер
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})
	//отключаем лог, если не установлен соответствующий флаг
	if !*config.LogToStdout {
		logrus.SetOutput(io.Discard)
	}

	logrus.Info("Service init stated")

	// создаем мастер-воркер, из которого будем обрабатывать таски. сделано так, чтобы можно было подсунуть другие
	// обработчики, реализующие task.WorkerInterface
	worker := taskWorker.NewWorker(config.GetWorkerTimeout())
	router := taskRouter.NewRouter(worker, *config.WorkerLimit)
	router.Run()

	// создаем эмиттер тасок, передаем ему входящий канал роутера и, если надо, *time.Duration для задержки между тасками
	creator := taskEmitter.NewEmitter(router.GetInputChannel(), config.GetTaskTimeout())
	creator.EmitTasks()

	logrus.Info("All systems online")

	time.Sleep(time.Second * time.Duration(*config.RunTime))

	logrus.Warn("Starting shutdown sequence")

	creator.Quit()
	router.Quit()

	if *config.PrintResults {
		router.PrintErrors()
		router.PrintFinished()
	}
}
