package collections

type QueueNode[T any] struct {
	Value *T
	Next  *QueueNode[T]
}

type Queue[T any] struct {
	head *QueueNode[T]
	tail *QueueNode[T]
	len  int
}

func (q *Queue[T]) Clear() {
	q.head = nil
	q.tail = nil
	q.len = 0
}

func (q *Queue[T]) Enqueue(value *T) {
	new_node := new(QueueNode[T])
	new_node.Value = value
	new_node.Next = q.tail
	q.tail = new_node

	if q.head == nil {
		q.head = q.tail
	}

	q.len++
	return
}

func (q *Queue[T]) Dequeue() (*T, error) {
	if q.len == 0 {
		return nil, nil
	}

	var node *QueueNode[T] = q.tail

	q.tail = node.Next

	if q.head == node {
		q.head = nil
	}

	q.len--

	return node.Value, nil
}

func (q *Queue[T]) Peek() *T {
	if q.len == 0 {
		return nil
	}

	return q.tail.Value
}

func (q *Queue[T]) Len() int {
	return q.len
}

func (q *Queue[T]) IsEmpty() bool {
	return q.len == 0
}
