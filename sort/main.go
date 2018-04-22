package main

import (
	"fmt"
)

/*****************************************************************************
1、简介
	SORT包中实现了3种基本的排序算法:插入排序、快排和堆排序.和其他语言中一样这三种方式都是不公开的
他们只在SORT包内部使用所以用户在使用SORT包进行排序时无需考虑使用那种排序方式.
SORT.INTERFACE定义的三个方法:
1、获取数据集合长度的LEN()方法
2、比较两个元素大小的LESS()方法
3、交换两个元素位置的SWAP()方法
就可以顺利对数据集合进行排序
SORT包会根据实际数据自动选择高效的排序算法



任何实现了SORT.INTERFACE的类型(一般为集合)均可使用该包中的方法进行排序
这些方法要求集合内列出元素的索引为整数

这个包中还有很多方法.这个包实现了很多方法比如排序反转、二分搜索
排序通过QUICKSORT()这个方法来控制该调用快排还是堆排
*****************************************************************************/

type Interface interface {
	Len() int
	Less(i, j int) bool
	Swap(i, j int)
}

// 1、源码中的例子:
type Person struct {
	Age int
}

type ByAge []Person

func (a ByAge) Len() int {
	return len(a)
}
func (a ByAge) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a ByAge) Less(i, j int) bool {
	return a[i].Age < a[j].Age
}

// 2、SORT方法

func Sort(data Interface) { // SORT包只提供了这一个公开使用的排序方法

	n := data.Len() // 如果元素深度达到2*CEIL(LG(N+1))则选用堆排序

	maxDepth := 0

	for i := n; i > 0; i >>= 1 {
		maxDepth++
	}

	maxDepth *= 2 // 2倍深度

	quickSort(data, 0, n, maxDepth)

}

func quickSort(data Interface, a, b, maxDepth int) { // 快速排序:它这里会自动选择是用堆排序还是插入排序还是快速排序

	// 第一个参数表示起始索引值
	// 第二个参数表示元素个数
	// 第三个参数表示2倍深度

	for b-a > 12 { // 如果切片元素少于十二个则使用希尔插入法

		if maxDepth == 0 {

			heapSort(data, a, b) // 堆排序方法: a=0,b=n

			return

		}

		maxDepth--

		mlo, mhi := doPivot(data, a, b)

		if mlo-a < b-mhi {
			quickSort(data, a, mlo, maxDepth)
			a = mhi
		} else {
			quickSort(data, mhi, b, maxDepth)
			b = mlo
		}

	}

	if b-a > 1 { // 希尔插入法

		for i := a + 6; i < b; i++ { // 对半开(元素最多不超过12个)
			if data.Less(i, i-6) {
				data.Swap(i, i-6)
			}
		}

		insertionSort(data, a, b)

	}
}
func insertionSort(data Interface, a, b int) { // 12个元素以内--插入排序

	for i := a + 1; i < b; i++ {

		for j := i; j > a && data.Less(j, j-1); j-- {
			data.Swap(j, j-1)
		}

	}
}

func heapSort(data Interface, a, b int) { // 堆排序

	first := a

	lo := 0
	hi := b - a

	for i := (hi - 1) / 2; i >= 0; i-- { // 构建堆结构最大的元素的顶部就是构建大根堆
		siftDown(data, i, hi, first)
	}

	//把first插入到data的end结尾
	for i := hi - 1; i >= 0; i-- {
		data.Swap(first, first+i)    //数据交换
		siftDown(data, lo, i, first) //堆重新筛选
	}
}

func siftDown(data Interface, lo, hi, first int) {
	//hi为数组的长度
	//这里有一种做法是把跟元素给取到存下来，但是为了方法更抽象，省掉了这部，取而代之的是在swap的时候进行相互交换
	root := lo //根元素的下标
	for {
		child := 2*root + 1 //左叶子结点下标
		//控制for循环介绍，这种写法更简洁，可以查看我写的堆排序的文章
		if child >= hi {
			break
		}
		//防止数组下标越界，判断左孩子和右孩子那个大
		if child+1 < hi && data.Less(first+child, first+child+1) {
			child++
		}
		//判断最大的孩子和根元素之间的关系
		if !data.Less(first+root, first+child) {
			return
		}
		//如果上面都 满足，则进行数据交换
		data.Swap(first+root, first+child)
		root = child
	}
}
func doPivot(data Interface, lo, hi int) (midlo, midhi int) {

	m := int(uint(lo+hi) >> 1)

	if hi-lo > 40 {

		s := (hi - lo) / 8

		medianOfThree(data, lo, lo+s, lo+2*s)
		medianOfThree(data, m, m-s, m+s)
		medianOfThree(data, hi-1, hi-1-s, hi-1-2*s)

	}

	medianOfThree(data, lo, m, hi-1)

	pivot := lo

	a, c := lo+1, hi-1

	for ; a < c && data.Less(a, pivot); a++ {
	}

	b := a

	for {

		for ; b < c && !data.Less(pivot, b); b++ {
		}

		for ; b < c && data.Less(pivot, c-1); c-- {
		}

		if b >= c {
			break
		}

		data.Swap(b, c-1)

		b++

		c--

	}

	protect := hi-c < 5

	if !protect && hi-c < (hi-lo)/4 {

		dups := 0

		if !data.Less(pivot, hi-1) {
			data.Swap(c, hi-1)
			c++
			dups++
		}

		if !data.Less(b-1, pivot) {
			b--
			dups++
		}

		if !data.Less(m, pivot) {
			data.Swap(m, b-1)
			b--
			dups++
		}

		protect = dups > 1
	}

	if protect {

		for {

			for ; a < b && !data.Less(b-1, pivot); b-- {
			}

			for ; a < b && data.Less(a, pivot); a++ {
			}

			if a >= b {
				break
			}

			data.Swap(a, b-1)

			a++
			b--
		}
	}

	data.Swap(pivot, b-1)

	return b - 1, c
}
func medianOfThree(data Interface, m1, m0, m2 int) { // 三个数之间的排序(参数为索引)

	// 按1、0、2排序

	if data.Less(m1, m0) {
		data.Swap(m1, m0)
	}

	if data.Less(m2, m1) {

		data.Swap(m2, m1)

		if data.Less(m1, m0) {
			data.Swap(m1, m0)
		}

	}

}

func main() {

	people := []Person{
		{7},
		{48},
		{13},
		{255},
		{66},
		{3},
		{67},
		{26},
		{31},
		{78},
		{17},
		{23},
		{321},
		{42},
		{683},
		{22},
	}

	fmt.Println(people)
	Sort(ByAge(people)) // 此处调用了SORT包中的SORT()方法
	fmt.Println(people)

}
