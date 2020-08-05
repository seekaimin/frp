package set

import (
	"container/list"
	"utils/common"
)

type setCache struct {
	value interface{}
	mark  bool
}

//Set 队列
type Set struct {
	items *list.List
	lock  commontools.MyLock
}

//New 创建一个新的队列
func New() *Set {
	n := new(Set)
	n.lock = commontools.MyLock{}
	n.items = list.New()
	return n
}

//Count 获取队列长度
func (p *Set) Count() int {
	return p.items.Len()
}

//Init 初始化
func (p *Set) Init() {
	p.lock.Lock(func() {
		p.items.Init()
	})
}

//First 取出(并移除)第一个
func (p *Set) First() interface{} {
	return p.lock.LockReturn(func() interface{} {
		item := p.items.Front()
		if item != nil {
			p.items.Remove(item)
			return item.Value
		}
		return nil
	})
}

//Last 取出(并移除)第最后一个
func (p *Set) Last() interface{} {
	return p.lock.LockReturn(func() interface{} {
		item := p.items.Back()
		if item != nil {
			p.items.Remove(item)
			return item.Value
		}
		return nil
	})
}

//PushFront 在最前面添加
func (p *Set) PushFront(val interface{}) {
	p.lock.Lock(func() {
		p.items.PushFront(val)
	})
}

//Push 在最后添加
func (p *Set) Push(val interface{}) {
	p.lock.Lock(func() {
		p.items.PushBack(val)
	})
}

//是否存在元素
func (p *Set) Has(val interface{}) bool {
	flag := false
	if val == nil {
		return flag
	}
	p.lock.Lock(func() {
		item := p.items.Front()
		for {
			if item == nil {
				break
			}
			if item.Value == val {
				flag = true
				break
			}
			item = item.Next()
		}
	})
	return flag
}

//移除
func (p *Set) Remove(val interface{}) {
	if val == nil {
		return
	}
	p.lock.Lock(func() {
		item := p.items.Front()
		for {
			if item == nil {
				break
			}
			if item.Value == val {
				p.items.Remove(item)
				break
			}
			item = item.Next()
		}
	})
}

//遍历
func (p *Set) Each(fun func(index int, item interface{})) {
	p.lock.Lock(func() {
		var index int
		if p.Count() > 0 {
			item := p.items.Front()
			for item != nil {
				fun(index, item.Value)
				item = item.Next()
				index++
			}
		}
	})
}
func (p *Set) EachLoop(fun func(index int, item interface{}) bool) {
	p.lock.Lock(func() {
		var index int
		if p.Count() > 0 {
			item := p.items.Front()
			for item != nil {
				if fun(index, item.Value) {
					break
				}
				item = item.Next()
				index++
			}
		}
	})
}
