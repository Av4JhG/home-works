package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	// Объявляем служебные переменные.
	done := make(chan interface{})
	taskChan := make(chan Task)
	resultChan := make(chan error, n)

	// Создаем Wait Group и n горутин.
	wg := &sync.WaitGroup{}
	for i := 0; i < n; i++ {
		wg.Add(1)
		go worker(wg, done, resultChan, taskChan)
	}

	// Запуск таски и определение превышения лимита ошибок.
	isErrorsLimitExceeded := runTasks(tasks, resultChan, taskChan, m)

	// Закрываем каналы и ждем завершения всех горутин.
	close(done)
	close(taskChan)
	wg.Wait()

	// Возвращаем при превышении лимита ошибок ошибку, или nil.
	if isErrorsLimitExceeded {
		return ErrErrorsLimitExceeded
	}
	return nil
}

// Функция обработки тасок.
func worker(wg *sync.WaitGroup, done <-chan interface{}, resultChan chan<- error, taskChan <-chan Task) {
	defer wg.Done()

	for task := range taskChan {
		result := task()
		select {
		case <-done:
			return
		case resultChan <- result:
		}
	}
}

// Функция отправки в канал тасок и отлов ошибок.
func runTasks(tasks []Task, resultChan <-chan error, taskChan chan<- Task, m int) bool {
	errorsCount := 0

	for i := 0; i < len(tasks); {
		task := tasks[i]

		select {
		case err := <-resultChan:
			if m > 0 && err != nil {
				errorsCount++
				if errorsCount == m {
					return true
				}
			}
		case taskChan <- task:
			i++
		}
	}

	return false
}
