package mailbox

type item struct {
	contents interface{}
	priority int
	index int
}

type priorityQueue []*item

func (q priorityQueue) Len() int {
	return len(q)
}

func (q priorityQueue) Less(i, j int) bool {
	return q[i].priority > q[j].priority
}

func (q priorityQueue) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
	q[i].index = i
	q[j].index = j
}

func (q *priorityQueue) Push(x interface{}) {
	n := len(*q)
	item := x.(*item)
	item.index = n
	*q = append(*q, item)
}

func (q *priorityQueue) Pop() interface{} {
	old := *q
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*q = old[0:n-1]
	return item
}
