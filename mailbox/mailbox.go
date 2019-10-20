package mailbox

import (
	"container/heap"
	"context"
)

type MailBox struct {
	PeerMailBox

	queue priorityQueue
	mode  QMode
	size  int

	putCh   chan *Message
	getCh   chan ResultCh
	waiting []ResultCh
}

type PeerMailBox interface {
	Put(message Message)
}

type Message struct {
	Contents interface{}
	Priority int
}

// The MailBox uses it's internal priorityQueue in two different ways: as an
// array or as a heap. This is because the `container.heap` implementation
// doesn't initialize properly unless the heap has more than 2 items.
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

		putCh:   make(chan *Message),
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
				resultCh <- item.Contents
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
				resultCh <- m.queue[0].Contents
				m.queue = m.queue[1:]
				break
			}

			if m.mode == Heap {
				message := heap.Pop(&m.queue).(*Message)
				resultCh <- message.Contents
				break
			}
		}
	}
}

func (m MailBox) Put(message Message) {
	m.putCh <- &message
}

func (m MailBox) Get(ctx context.Context) interface{} {
	resultCh := make(ResultCh)
	m.getCh <- resultCh
	select {
	case contents := <-resultCh:
		return contents
	case <-ctx.Done():
		return nil
	}
}
