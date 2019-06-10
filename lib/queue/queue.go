package queue

type Item struct {
	Key   []byte
	Value []byte
}

type Queue interface {
	Enqueue(clickId string, value []byte) error
	Dequeue() (*Item, error)
}
