package main

import (
	"sync"
	"sync/atomic"
	"unsafe"
	"fmt"
	"time"
)

type Node struct{
	num int
	next *Node
}

type Queue struct {
	head *Node
	tail *Node
	hmutex sync.Mutex
	tmutex sync.Mutex
}

func (q *Queue)len() int { 
	length:=0
	s:=q.head.next

	for(s!=nil){
		s=s.next
		length+=1
	}
	return length
}

func NewQueue() *Queue {
	n := &Node{} // dummy 节点
	return &Queue{
		head: n,
		tail: n,
	}
}

func (q *Queue)Enqueue(num int)(bool) {
	node:=&Node{num:num,next:nil}
	q.tmutex.Lock()
	defer q.tmutex.Unlock()
	// tail受保护但是tail.next不一定
	// atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&q.tail.next)),(unsafe.Pointer(node)))
	q.tail.next = node
	q.tail=node
	return true
}

func (q *Queue)Dequeue()(*Node){
	q.tmutex.Lock()
	defer q.tmutex.Unlock()
	next := load(q.head.next)
	if(next == nil){
		return nil
	}
	q.head=next
	return next
}

func load(node *Node)(*Node){
	return (*Node)(atomic.LoadPointer(((*unsafe.Pointer)(unsafe.Pointer(&node)))));
}

func main() {
	start_time:=time.Now()
	defer func(){ 
		fmt.Println(time.Since(start_time))
	}()
	counts:=50000
	n := &Node{num:0,next:nil}
	q := Queue{head:n,tail:n}
	var wg sync.WaitGroup
	for i := 0; i <800;i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := i*counts; j <(i+1)*counts; j++{
				q.Enqueue(j)
				// fmt.Println("over goroutine",j)
			}
			for j := i*counts; j <(i+1)*counts; j++{
				q.Dequeue()
				// fmt.Println("over goroutine",j)
			}			
		}(i)
	}
	wg.Wait()
	// fmt.Println(q.len())
}
