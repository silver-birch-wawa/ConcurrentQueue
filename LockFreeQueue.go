package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"unsafe"
	"time"
)
type Node struct{
	num int
	next *Node
}

type Queue struct {
	head *Node
	tail *Node
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
func (q *Queue)Enqueue(num int){
	newnode:=&Node{num:num,next:nil}
	for{
		tail:=(*Node)(atomic.LoadPointer(((*unsafe.Pointer)(unsafe.Pointer(&q.tail)))))
		next:=tail.next
		if next==nil{
			if(atomic.CompareAndSwapPointer(((*unsafe.Pointer)(unsafe.Pointer(&tail.next))), unsafe.Pointer(nil),unsafe.Pointer(newnode))){
				atomic.CompareAndSwapPointer(((*unsafe.Pointer)(unsafe.Pointer(&q.tail))), unsafe.Pointer(tail),unsafe.Pointer(tail.next))
				break
			}
		}
		atomic.CompareAndSwapPointer(((*unsafe.Pointer)(unsafe.Pointer(&q.tail))), unsafe.Pointer(tail),unsafe.Pointer(tail.next))
	}
}
func (q *Queue)Dequeue()(*Node) {
	for{
        head := (*Node)(atomic.LoadPointer(((*unsafe.Pointer)(unsafe.Pointer(&q.head)))))
        tail := (*Node)(atomic.LoadPointer(((*unsafe.Pointer)(unsafe.Pointer(&q.tail)))))
        next := (*Node)(atomic.LoadPointer(((*unsafe.Pointer)(unsafe.Pointer(&head.next)))))
		if(next==nil){
			return nil
		}
		if(cas(&q.head,head,next)){
			return next
		}
		if(head==tail){
			// next==nil但是head却==tail说明有人在尾部插入,那么就帮插入圆梦
			atomic.CompareAndSwapPointer(((*unsafe.Pointer)(unsafe.Pointer(&q.tail))), unsafe.Pointer(tail),unsafe.Pointer(tail.next))
		}
	}
}
func cas(node1 **Node, node2 *Node, node3 *Node) (ok bool) {
    return atomic.CompareAndSwapPointer(
		// 第一个是unsafe pointer的地址,第二个是unsafe pointer ,第三个也是unsafe pointer
		((*unsafe.Pointer)(unsafe.Pointer(node1))), unsafe.Pointer(node2),unsafe.Pointer(node3))
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