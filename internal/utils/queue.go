package utils

type Queue[T any] struct {
	arr []T
}

func (q *Queue[T]) Push(val T) {
	q.arr = append(q.arr, val)
}

func (q *Queue[T]) Pop() (el T, success bool) {
	if len(q.arr) == 0 {
		return
	}
	el = q.arr[0]
	q.arr = q.arr[1:]
	return el, true
}

func (q *Queue[T]) Len() int {
	return len(q.arr)
}
