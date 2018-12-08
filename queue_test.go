package conc_queue

import (
	"sync"
	"testing"
)

const limit = 10000

func TestQueue(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(2)
	q := NewQueue(10)
	go func() {
		defer wg.Done()
		for i := 1; i <= limit; i++ {
			q.Push(i)
		}
	}()
	var val int
	go func(valPtr *int) {
		defer wg.Done()
		for i := 1; i <= limit; i++ {
			*valPtr += q.Pop().(int)
		}
	}(&val)
	wg.Wait()
	exp := (limit + 1) * (limit >> 1)
	if val != exp {
		t.Errorf("Expected conuter to be %d. Actual is %d", exp, val)
	}
}
