package mailbox

type priorityQueue []*qItem

type qItem struct {
	contents interface{}
	priority int
	index    int
}

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
	q[i].index = i
	q[j].index = j
}

func (q *priorityQueue) Push(x interface{} /* qItem */) {
	n := len(*q)
	item := x.(*qItem)
	item.index = n
	*q = append(*q, item)
}

// Returns a *qItem.
func (q *priorityQueue) Pop() interface{} /* qItem */ {
	oldQ := *q
	size := len(oldQ)
	item := oldQ[size-1]
	oldQ[size-1] = nil // avoid memory leak
	item.index = -1    // for safety
	*q = oldQ[0 : size-1]
	return item
}
