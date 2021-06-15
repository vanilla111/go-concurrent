package list

import (
	"fmt"
	"sync"
	"sync/atomic"
	"unsafe"
)

type IntList struct {
	head   *Node
	length int64
	mu     sync.Mutex
}

type Node struct {
	value  int
	marked bool
	next   *Node
	mu     sync.Mutex
}

func newNode(value int) *Node {
	return &Node{value: value}
}

func NewInt() *IntList {
	// init with head
	return &IntList{head: newNode(0)}
}

func (l *IntList) Contains(value int) bool {
	B := loadNext(l.head)
	for B != nil {
		if B.value == value {
			if B.marked {
				return false
			}
			return true
		}
		B = loadNext(B)
	}
	return false
}

func (l *IntList) Insert(value int) bool {
	node := newNode(value)
	RESTART:
		var A = l.head
		var B = loadNext(A)
		for B != nil && B.value <= value {
			if B.value == value {
				goto FOUND
			}
			A = B
			B = loadNext(A)
		}


	FOUND:
		A.mu.Lock()
		if loadNext(A) != B || A.marked {
			A.mu.Unlock()
			goto RESTART
		}
		if B != nil && B.value == value {
			A.mu.Unlock()
			return false
		}
		node.next = B
		storeNext(A, node)
		A.mu.Unlock()
		l.head.mu.Lock()
		defer l.head.mu.Unlock()
		l.length++
		return true
}

func (l *IntList) Delete(value int) bool {
	FIND:
	A := l.head
	B := loadNext(A)
	for B != nil && B.value != value {
		A = B
		B = loadNext(A)
	}
	if B == nil {
		return false
	}
	B.mu.Lock()
	if B.marked {
		B.mu.Unlock()
		goto FIND
	}
	A.mu.Lock()
	if A.next != B || A.marked {
		A.mu.Unlock()
		B.mu.Unlock()
		goto FIND
	}
	B.marked = true
	storeNext(A, B.next)
	A.mu.Unlock()
	B.mu.Unlock()
	l.head.mu.Lock()
	defer l.head.mu.Unlock()
	l.length--
	return true
}

func (l *IntList) Range(f func(value int) bool) {
	B := loadNext(l.head)
	for B != nil {
		if !f(B.value) {
			break
		}
		B = loadNext(B)
	}
}

func (l *IntList) Len() int {
	l.head.mu.Lock()
	defer l.head.mu.Unlock()
	return int(l.length)
}

func (l *IntList) Print() {
	var p = l.head
	for p.next != nil {
		if p.next.value <= p.value {
			panic("List不符合语义")
		}
		fmt.Print(p.next.value, "-> ")
		p = p.next
	}
	fmt.Println(nil)
}

func loadNext(node *Node) *Node {
	return (*Node)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&node.next))))
}

func storeNext(pre, new *Node) {
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&pre.next)), unsafe.Pointer(new))
}
