package mailbox

import (
	"container/heap"
	"context"
)

type MailBox struct {
	queue *priorityQueue
	size  int
	putCh chan *item
	getCh chan chan interface{}
	resultCh chan interface{}
}

func NewMailBox(size int) MailBox {
	queue := make(priorityQueue, size)
	heap.Init(&queue)
	box := MailBox{
		queue: &queue,
		size:  size,
		putCh: make(chan *item),
		getCh: make(chan chan interface{}),
	}
	go box.doOps()
	return box
}

func (m MailBox) doOps() {
	for {
		select {
		case item := <-m.putCh:
			if m.resultCh != nil {
				m.resultCh <- item.contents
				m.resultCh = nil
				break
			}

			if m.queue.Len() < m.size {
				heap.Push(m.queue, item)
			}

		case resultCh := <-m.getCh:
			if m.queue.Len() == 0 {
				if m.resultCh != nil {
					panic("double get")
				}
				m.resultCh = resultCh
				break
			}

			item := heap.Pop(m.queue).(*item)
			resultCh <- item.contents
		}
	}
}

func (m MailBox) Put(message interface{}, priority int) {
	m.putCh <- &item{
		contents: message,
		priority: priority,
	}
}

func (m MailBox) Get(ctx context.Context) interface{} {
	resultCh := make(chan interface{})
	m.getCh <- resultCh
	select {
	case message := <-resultCh:
		return message
	case <-ctx.Done():
		return nil
	}
}
