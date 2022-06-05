package internal

import (
	"log"
)

const messageQueueSize = 100

type Message struct {
	Payload []byte
}

type Handler interface {
	Handle(message *Message) error
}

func NewQueue() *Queue {
	return &Queue{ch: make(chan *Message, messageQueueSize)}
}

type Queue struct {
	ch chan *Message
}

func (q *Queue) IsEmpty() bool {
	return len(q.ch) == 0
}

func (q *Queue) Close() {
	close(q.ch)
}

func (q *Queue) Send(message *Message) {
	q.ch <- message
}

func (q *Queue) Listen(handler Handler) {
	for msg := range q.ch {
		if err := handler.Handle(msg); err != nil {
			log.Println(err)
		}
	}
}
