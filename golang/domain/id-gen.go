package domain

import "sync/atomic"

var currentId atomic.Int32

// в оригинале мы генерили ID таска из времени, но у нас нет гарантий,
// что мы не сгенерим 2 таска с одним и тем же ID

func getNextTaskId() int {
	currentId.Add(1)
	return int(currentId.Load())
}
