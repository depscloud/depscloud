package eventlp

// T is a definition for a generic value stored in the linked list
type T interface{}

type node struct {
	prev *node
	next *node
	val  T
}

// LinkedList implements a doubly linked list making it easy to treat as a
// stack or as a queue. This shares a similar interface to c++'s vector object.
type LinkedList struct {
	head *node
	tail *node
	size int
}

// PushFront prepends an element to the start of the list.
func (ll *LinkedList) PushFront(t T) {
	n := &node{
		val: t,
	}

	if ll.head == nil {
		ll.head = n
		ll.tail = n
	} else {
		n.next = ll.head
		ll.head.prev = n
		ll.head = n
	}

	ll.size++
}

// PopFront removes the element from the start of the list.
func (ll *LinkedList) PopFront() T {
	if ll.head == nil {
		return nil
	}

	val := ll.head.val
	if ll.head == ll.tail {
		ll.head = nil
		ll.tail = nil
	} else {
		head := ll.head.next
		head.prev = nil
		ll.head.next = nil
		ll.head = head
	}

	ll.size--
	return val
}

// PushBack appends an element to the end of the list.
func (ll *LinkedList) PushBack(t T) {
	n := &node{
		val: t,
	}

	if ll.tail == nil {
		ll.head = n
		ll.tail = n
	} else {
		n.prev = ll.tail
		ll.tail.next = n
		ll.tail = n
	}

	ll.size++
}

// PopBack removes the element at the end of the list.
func (ll *LinkedList) PopBack() T {
	if ll.tail == nil {
		return nil
	}

	val := ll.tail.val
	if ll.head == ll.tail {
		ll.head = nil
		ll.tail = nil
	} else {
		tail := ll.tail.prev
		tail.next = nil
		ll.tail.prev = nil
		ll.tail = tail
	}

	ll.size--
	return val
}

func (ll *LinkedList) Size() int {
	return ll.size
}
