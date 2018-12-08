package conc_queue

import (
	"sync"
)

type Queue struct {
	head, tail, size, capacity uint
	arr                        []interface{}
	fullCond, emptyCond        *sync.Cond
	mutex                      *sync.Mutex
}

func NewQueue(capacity uint) *Queue {
	return &Queue{
		capacity:  capacity,
		arr:       make([]interface{}, capacity),
		fullCond:  sync.NewCond(&sync.Mutex{}),
		emptyCond: sync.NewCond(&sync.Mutex{}),
		mutex:     &sync.Mutex{},
	}
}

func (this *Queue) Push(val interface{}) {
	this.ensureNotFull()
	this.push(val)
	this.emptyCond.Signal()
}

func (this *Queue) Pop() interface{} {
	this.ensureNotEmpty()
	result := this.pop()
	this.fullCond.Signal()
	return result
}

func (this *Queue) ensureNotFull() {
	this.fullCond.L.Lock()
	defer this.fullCond.L.Unlock()
	for {
		this.mutex.Lock()
		if this.size == this.capacity {
			this.mutex.Unlock()
			this.fullCond.Wait()
		} else {
			this.mutex.Unlock()
			break
		}
	}
}

func (this *Queue) ensureNotEmpty() {
	this.emptyCond.L.Lock()
	defer this.emptyCond.L.Unlock()
	for {
		this.mutex.Lock()
		if this.size == 0 {
			this.mutex.Unlock()
			this.emptyCond.Wait()
		} else {
			this.mutex.Unlock()
			break
		}
	}
}

func (this *Queue) push(val interface{}) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	this.arr[this.tail] = val
	//this.tail = (this.tail + 1) & (this.capacity - 1)
	this.tail = (this.tail + 1) % this.capacity
	this.size++
}

func (this *Queue) pop() interface{} {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	result := this.arr[this.head]
	//this.head = (this.head + 1) & (this.capacity - 1)
	this.head = (this.head + 1) % this.capacity
	this.size--
	return result
}
