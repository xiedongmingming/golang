package main

import (
	"errors"
	"reflect"
	"sync"
)

var errBadChannel = errors.New("event: Subscribe argument does not have sendable channel type")

const firstSubSendCase = 1

type feedTypeError struct {
	got, want reflect.Type
	op        string
}

func (e feedTypeError) Error() string {
	return "event: wrong type in " + e.op + " got " + e.got.String() + ", want " + e.want.String()
}

// ***************************************************************************
type Feed struct { //
	once     sync.Once     // 确保初始化只执行一次
	sendLock chan struct{} // 作为锁使用

	removeSub chan interface{} // 待移除的订阅--将准备移除的订阅发送到该通道中

	sendCases caseList // []reflect.SelectCase -- 用于接收
	inbox     caseList // []reflect.SelectCase -- 用于发送(订阅通道放在此处)

	mu    sync.Mutex   //
	etype reflect.Type // 表示通道的数据类型

	closed bool
}

func (f *Feed) init() { // 只会执行一次

	f.removeSub = make(chan interface{})

	f.sendLock = make(chan struct{}, 1)
	f.sendLock <- struct{}{}

	f.sendCases = caseList{{Chan: reflect.ValueOf(f.removeSub), Dir: reflect.SelectRecv}}

}
func (f *Feed) Subscribe(channel interface{}) Subscription { // 订阅的意思是将事件发送到参数表示的通道中

	f.once.Do(f.init)

	chanval := reflect.ValueOf(channel)

	chantyp := chanval.Type()

	if chantyp.Kind() != reflect.Chan || chantyp.ChanDir()&reflect.SendDir == 0 { // 确保是通道并且可以发送
		panic(errBadChannel)
	}

	sub := &feedSub{feed: f, channel: chanval, err: make(chan error, 1)}

	f.mu.Lock()
	defer f.mu.Unlock()

	if !f.typecheck(chantyp.Elem()) {
		panic(feedTypeError{op: "subscribe", got: chantyp, want: reflect.ChanOf(reflect.SendDir, f.etype)})
	}

	cas := reflect.SelectCase{Dir: reflect.SelectSend, Chan: chanval}

	f.inbox = append(f.inbox, cas)

	return sub
}
func (f *Feed) typecheck(typ reflect.Type) bool {

	if f.etype == nil {

		f.etype = typ

		return true
	}

	return f.etype == typ
}
func (f *Feed) remove(sub *feedSub) { // 将一个代表订阅通道的包装移除

	ch := sub.channel.Interface() // 获取外部通道的值

	f.mu.Lock()

	index := f.inbox.find(ch)

	if index != -1 {

		f.inbox = f.inbox.delete(index)

		f.mu.Unlock()

		return
	}

	f.mu.Unlock()

	// 当上面没找到时????
	select {
	case f.removeSub <- ch:
	case <-f.sendLock:
		f.sendCases = f.sendCases.delete(f.sendCases.find(ch))
		f.sendLock <- struct{}{}
	}
}
func (f *Feed) Send(value interface{}) (nsent int) {

	f.once.Do(f.init)

	<-f.sendLock // 发送锁

	f.mu.Lock()
	f.sendCases = append(f.sendCases, f.inbox...)
	f.inbox = nil
	f.mu.Unlock()

	rvalue := reflect.ValueOf(value)

	if !f.typecheck(rvalue.Type()) {

		f.sendLock <- struct{}{}

		panic(feedTypeError{op: "send", got: rvalue.Type(), want: f.etype})
	}

	for i := firstSubSendCase; i < len(f.sendCases); i++ {
		f.sendCases[i].Send = rvalue
	}

	cases := f.sendCases

	for {

		for i := firstSubSendCase; i < len(cases); i++ {
			if cases[i].Chan.TrySend(rvalue) {
				nsent++
				cases = cases.deactivate(i)
				i--
			}
		}
		if len(cases) == firstSubSendCase {
			break
		}

		chosen, recv, _ := reflect.Select(cases)
		if chosen == 0 {
			index := f.sendCases.find(recv.Interface())
			f.sendCases = f.sendCases.delete(index)
			if index >= 0 && index < len(cases) {
				cases = f.sendCases[:len(cases)-1]
			}
		} else {
			cases = cases.deactivate(chosen)
			nsent++
		}
	}

	for i := firstSubSendCase; i < len(f.sendCases); i++ {
		f.sendCases[i].Send = reflect.Value{}
	}

	f.sendLock <- struct{}{}

	return nsent
}

// ***************************************************************************
type feedSub struct { // 表示包装外部提供的一个订阅通道--实现接口: Subscription
	feed    *Feed         // 宿主
	channel reflect.Value // 外部的通道
	errOnce sync.Once     //
	err     chan error    //
}

func (sub *feedSub) Unsubscribe() {
	sub.errOnce.Do(func() {
		sub.feed.remove(sub)
		close(sub.err)
	})
}
func (sub *feedSub) Err() <-chan error {
	return sub.err
}

// ***************************************************************************
type caseList []reflect.SelectCase

func (cs caseList) find(channel interface{}) int {

	for i, cas := range cs {

		if cas.Chan.Interface() == channel {
			return i
		}
	}

	return -1
}
func (cs caseList) delete(index int) caseList {
	return append(cs[:index], cs[index+1:]...)
}
func (cs caseList) deactivate(index int) caseList {
	last := len(cs) - 1
	cs[index], cs[last] = cs[last], cs[index]
	return cs[:last]
}

// ***************************************************************************
// func (cs caseList) String() string {
//     s := "["
//     for i, cas := range cs {
//             if i != 0 {
//                     s += ", "
//             }
//             switch cas.Dir {
//             case reflect.SelectSend:
//                     s += fmt.Sprintf("%v<-", cas.Chan.Interface())
//             case reflect.SelectRecv:
//                     s += fmt.Sprintf("<-%v", cas.Chan.Interface())
//             }
//     }
//     return s + "]"
// }
// ***************************************************************************
type WalletEventType int

const ( // 三种类型
	WalletArrived WalletEventType = iota
	WalletOpened
	WalletDropped
)

type WalletEvent struct { // 表示钱包事件
	Wallet Wallet          // 包含的钱包
	Kind   WalletEventType // 钱包事件类型
}

func main() {

	var feed Feed

	sink := make(chan WalletEvent)

	feed.Subscribe(sink)

	feed.Send(event)
}
