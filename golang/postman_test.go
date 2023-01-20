package main

import (
	"testing"
	"time"
)

//проверяем, что постман создает сообщения
//и за 3.5 секунд он создаст 3 сообщения
func TestPostmanCreate(t *testing.T) {
	postman := NewPostman(3500*time.Millisecond, 1*time.Second)
	tasks := postman.Create()
	time.Sleep(4 * time.Second)
	i := 0
	for range tasks {
		i++
	}
	if i != 3 {
		t.Fatalf("assertions don't eq %d !=%d", 2, i)
	}
}
