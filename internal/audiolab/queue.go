package audiolab

// Queue is a fixed-capacity event ring. Gameplay can enqueue effects without
// allocating a new slice or player per hit.
type Queue struct {
	values     []string
	head, size int
}

func NewQueue(capacity int) Queue { return Queue{values: make([]string, capacity)} }
func (q *Queue) Push(v string) {
	if len(q.values) == 0 {
		return
	}
	if q.size < len(q.values) {
		q.values[(q.head+q.size)%len(q.values)] = v
		q.size++
		return
	}
	q.values[q.head] = v
	q.head = (q.head + 1) % len(q.values)
}
func (q *Queue) Pop() (string, bool) {
	if q.size == 0 {
		return "", false
	}
	v := q.values[q.head]
	q.head = (q.head + 1) % len(q.values)
	q.size--
	return v, true
}
func (q *Queue) Len() int { return q.size }
