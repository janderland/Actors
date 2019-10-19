package mailbox

import (
	"container/heap"
	"context"
)

type MailBox struct {
	queue priorityQueue
	mode  QMode
	size  int

	putCh   chan *qItem
	getCh   chan ResultCh
	waiting []ResultCh
}

type QMode int

const (
	Array QMode = iota
	Heap
)

type ResultCh chan interface{}

func NewMailBox(size int) MailBox {
	box := MailBox{
		queue: nil,
		mode:  Array,
		size:  size,

		putCh:   make(chan *qItem),
		getCh:   make(chan ResultCh),
		waiting: nil,
	}

	go box.doSynchronization()

	return box
}

func (m MailBox) doSynchronization() {
	for {
		select {
		case item := <-m.putCh:
			if len(m.waiting) > 0 {
				var resultCh ResultCh
				resultCh, m.waiting = m.waiting[0], m.waiting[1:]
				resultCh <- item.contents
				break
			}

			if m.mode == Array {
				m.queue.Push(item)
				if m.queue.Len() > 1 {
					heap.Init(&m.queue)
					m.mode = Heap
				}
				break
			}

			if m.mode == Heap {
				if m.queue.Len() < m.size {
					heap.Push(&m.queue, item)
				}
				break
			}

		case resultCh := <-m.getCh:
			if m.queue.Len() == 0 {
				m.waiting = append(m.waiting, resultCh)
				break
			}

			if m.mode == Array {
				resultCh <- m.queue[0]
				m.queue = m.queue[1:]
				break
			}

			if m.mode == Heap {
				item := heap.Pop(&m.queue).(*qItem)
				resultCh <- item.contents
				break
			}
		}
	}
}

func (m MailBox) Put(message interface{}, priority int) {
	m.putCh <- &qItem{contents: message, priority: priority}
}

func (m MailBox) Get(ctx context.Context) interface{} {
	resultCh := make(ResultCh)
	m.getCh <- resultCh
	select {
	case message := <-resultCh:
		return message
	case <-ctx.Done():
		return nil
	}
}
