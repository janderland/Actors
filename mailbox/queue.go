package mailbox

type priorityQueue []*Message

func (q priorityQueue) Len() int {
	return len(q)
}

func (q priorityQueue) Less(i, j int) bool {
	// Use 'greater than' here because we want higher priorities to pop first.
	// The container.heap is designed to pop the lowest priority value first.
	return q[i].priority > q[j].priority
}

func (q priorityQueue) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
}

func (q *priorityQueue) Push(x interface{} /* Message */) {
	item := x.(*Message)
	*q = append(*q, item)
}

func (q *priorityQueue) Pop() interface{} /* Message */ {
	oldQ := *q
	size := len(oldQ)
	item := oldQ[size-1]
	oldQ[size-1] = nil // avoid memory leak
	*q = oldQ[0 : size-1]
	return item
}
