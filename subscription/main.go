package main

import (
	"errors"
	"fmt"
	"sync"
)

type Subscription interface {
	Err() <-chan error //
	Unsubscribe()      // 取消订阅
}

func NewSubscription(producer func(<-chan struct{}) error) Subscription { // 参数是一个函数

	s := &funcSub{unsub: make(chan struct{}), err: make(chan error, 1)}

	go func() { // 作为独立线程运行

		defer close(s.err)

		err := producer(s.unsub) // 会阻塞--处理通道

		s.mu.Lock()
		defer s.mu.Unlock()

		if !s.unsubscribed {
			if err != nil {
				s.err <- err
			}
			s.unsubscribed = true
		}
	}()

	return s // 直接返回(通过该对象向通道中发送数据就可以回调传递的函数了)
}

// ***************************************************************************
type funcSub struct { // 默认的订阅实现
	unsub        chan struct{} // 提供给外部处理的通道
	err          chan error    // 缓存为1--外部处理返回的错误结果
	mu           sync.Mutex    //
	unsubscribed bool          // 指示是否被取消订阅
}

func (s *funcSub) Unsubscribe() {

	s.mu.Lock()

	if s.unsubscribed {
		s.mu.Unlock()
		return
	}

	s.unsubscribed = true

	close(s.unsub) // 关闭外部处理通道

	s.mu.Unlock()

	<-s.err // 等待外部通道处理结束
}
func (s *funcSub) Err() <-chan error {
	return s.err
}

func Test(channel <-chan struct{}) error {

	fmt.Println("test ... ")

	return errors.New("test")
}

func main() {

	sub := NewSubscription(Test)
	
	(*funcSub)sub

	<-sub.Err()
}
