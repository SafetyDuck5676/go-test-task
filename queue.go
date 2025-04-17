package main

import (
	"errors"
	"sync"
	"time"
)

// Queue представляет одну очередь сообщений.
type Queue struct {
	messages []string
	maxSize  int
	mu       sync.Mutex
	waiters  []*waiter
}

// waiter представляет ожидающего получателя.
type waiter struct {
	ch chan string
}

// NewQueue создает новую очередь с максимальным размером maxSize.
func NewQueue(maxSize int) *Queue {
	return &Queue{
		messages: []string{},
		maxSize:  maxSize,
	}
}

// Add добавляет сообщение в очередь.
func (q *Queue) Add(message string) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.messages) >= q.maxSize {
		return errors.New("queue is full")
	}

	q.messages = append(q.messages, message)
	q.notifyWaiters()
	return nil
}

// Get извлекает сообщение из очереди по принципу FIFO.
func (q *Queue) Get(timeout time.Duration) (string, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.messages) > 0 {
		message := q.messages[0]
		q.messages = q.messages[1:]
		return message, nil
	}

	if timeout == 0 {
		return "", errors.New("no messages available")
	}

	// Создаем нового ожидающего получателя.
	waiter := &waiter{ch: make(chan string)}
	q.waiters = append(q.waiters, waiter)
	q.mu.Unlock()

	select {
	case message := <-waiter.ch:
		q.mu.Lock()
		return message, nil
	case <-time.After(timeout):
		q.mu.Lock()
		q.removeWaiter(waiter)
		return "", errors.New("timeout waiting for message")
	}
}

// notifyWaiters уведомляет всех ожидающих получателей о новом сообщении.
func (q *Queue) notifyWaiters() {
	for len(q.waiters) > 0 && len(q.messages) > 0 {
		// Переименуем переменную waiter в w, чтобы избежать конфликта имен.
		w := q.waiters[0]
		q.waiters = q.waiters[1:]
		message := q.messages[0]
		q.messages = q.messages[1:]
		go func(waiter *waiter, msg string) {
			waiter.ch <- msg
		}(w, message)
	}
}

// removeWaiter удаляет ожидающего получателя из списка.
func (q *Queue) removeWaiter(w *waiter) {
	for i, waiter := range q.waiters {
		if waiter == w {
			q.waiters = append(q.waiters[:i], q.waiters[i+1:]...)
			break
		}
	}
}

// QueueManager управляет всеми очередями.
type QueueManager struct {
	queues      map[string]*Queue
	maxQueues   int
	defaultSize int
	mu          sync.Mutex
}

// NewQueueManager создает новый менеджер очередей.
func NewQueueManager(maxQueues, defaultSize int) *QueueManager {
	return &QueueManager{
		queues:      make(map[string]*Queue),
		maxQueues:   maxQueues,
		defaultSize: defaultSize,
	}
}

// GetOrCreateQueue возвращает существующую очередь или создает новую.
func (qm *QueueManager) GetOrCreateQueue(name string) (*Queue, error) {
	qm.mu.Lock()
	defer qm.mu.Unlock()

	if len(qm.queues) >= qm.maxQueues {
		return nil, errors.New("maximum number of queues reached")
	}

	queue, exists := qm.queues[name]
	if !exists {
		queue = NewQueue(qm.defaultSize)
		qm.queues[name] = queue
	}
	return queue, nil
}
