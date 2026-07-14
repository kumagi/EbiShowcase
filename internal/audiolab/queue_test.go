package audiolab

import "testing"

func TestQueueReusesFixedCapacity(t *testing.T) {
	q := NewQueue(2)
	q.Push("a")
	q.Push("b")
	q.Push("c")
	if q.Len() != 2 {
		t.Fatal(q.Len())
	}
	a, _ := q.Pop()
	b, _ := q.Pop()
	if a != "b" || b != "c" {
		t.Fatalf("%q %q", a, b)
	}
}
