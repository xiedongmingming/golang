package main

import (
	"container/heap"
	"fmt"
	"sort"
)

// *****************************************************************************
type HeapInt []int

func (h HeapInt) Len() int {
	return len(h)
}
func (h HeapInt) Less(i, j int) bool {
	return h[i] < h[j]
}
func (h HeapInt) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}
func (h *HeapInt) Push(x interface{}) {
	*h = append(*h, x.(int))
}
func (h *HeapInt) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// *****************************************************************************
func main() {

	var slice1 *HeapInt

	slice1 = new(HeapInt)

	// heap.Init(&slice1)

	heap.Push(slice1, 23)
	heap.Push(slice1, 43)
	heap.Push(slice1, 1)
	heap.Push(slice1, 5)
	heap.Push(slice1, 65)
	heap.Push(slice1, 55)
	heap.Push(slice1, 88)
	heap.Push(slice1, 98)
	heap.Push(slice1, 8)
	heap.Push(slice1, 18)

	for i := 0; i < 10; i++ {
		fmt.Println((*slice1)[i])
	}
	fmt.Println("************************")

	//	for i := 0; i < 10; i++ {
	//		fmt.Println(heap.Pop(slice1))
	//	}
	fmt.Println("************************")

	sort.Sort(*slice1)

	for i := 0; i < 10; i++ {
		fmt.Println((*slice1)[i])
	}
}
