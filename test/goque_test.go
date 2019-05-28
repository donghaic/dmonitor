package test

import (
	"fmt"
	"github.com/beeker1121/goque"
	"testing"
	"time"
)

func TestEmbededQueue(t *testing.T) {
	q, _ := goque.OpenQueue("logs")
	defer q.Close()

	start := time.Now()
	for i := 1; i <= 10000; i++ {
		_, _ = q.EnqueueString(fmt.Sprintf("item %d", i))
	}
	fmt.Println(time.Since(start))
	println("queue len=", q.Length())
	start2 := time.Now()
	for {
		item, _ := q.Dequeue()
		if item != nil {
			println("data=", string(item.Value))
		} else {
			break
		}
	}
	fmt.Println(time.Since(start2))

}
