package eventlp_test

import (
	"testing"

	"github.com/depscloud/depscloud/internal/eventlp"

	"github.com/stretchr/testify/require"
)

func TestLinkedList_Queue(t *testing.T) {
	// ensures the linked list can behave like a queue
	queue := &eventlp.LinkedList{}

	queue.PushBack("a")
	queue.PushBack("b")
	queue.PushBack("c")

	require.Equal(t, 3, queue.Size())

	require.Equal(t, "a", queue.PopFront())
	require.Equal(t, "b", queue.PopFront())
	require.Equal(t, "c", queue.PopFront())
	require.Nil(t, queue.PopFront())

	require.Equal(t, 0, queue.Size())
}

func TestLinkedList_Stack(t *testing.T) {
	// ensures the linked list can behave like a stack
	stack := &eventlp.LinkedList{}

	stack.PushBack("a")
	stack.PushBack("b")
	stack.PushBack("c")

	require.Equal(t, 3, stack.Size())

	require.Equal(t, "c", stack.PopBack())
	require.Equal(t, "b", stack.PopBack())
	require.Equal(t, "a", stack.PopBack())
	require.Nil(t, stack.PopBack())

	require.Equal(t, 0, stack.Size())
}
