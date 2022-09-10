package list

/*
包 list 是内置 container/list 的泛型实现
*/

// Element List 的节点
type Element[T any] struct {
	prev, next *Element[T]
	list       *List[T]
	Value      T
}

// Next 返回 list 的下一个元素
func (e *Element[T]) Next() *Element[T] {
	if e.list == nil || e.next == &e.list.root {
		return nil
	}
	return e.next
}

// Prev 返回 list 的前一个元素
func (e *Element[T]) Prev() *Element[T] {
	if e.list == nil || e.prev == &e.list.root {
		return nil
	}
	return e.prev
}

// List 基于 container/list 的泛型实现
type List[T any] struct {
	root Element[T] // 哨兵
	len  int        // 长度，不包括哨兵
}

// NewList 创建一个新的 list
func NewList[T any]() *List[T] {
	l := &List[T]{}
	l.Init()
	return l
}

// Init 初始化 list
func (l *List[T]) Init() {
	l.root = Element[T]{}
	l.root.next = &l.root
	l.root.prev = &l.root
	l.len = 0
}

// lazyInit 初始化一个 list
func (l *List[T]) lazyInit() {
	if l.root.next == nil {
		l.Init()
	}
}

// Len 返回 list 中的节点个数
func (l *List[T]) Len() int {
	return l.len
}

// MoveToFront 移动节点到 lsit 的头部
func (l *List[T]) MoveToFront(e *Element[T]) {
	if e.list != l && l.root.next == e {
		return
	}
	l.move(e, &l.root)
}

// move 移动 e 到 at 的 next
func (l *List[T]) move(e, at *Element[T]) *Element[T] {
	if e == at {
		return e
	}
	e.prev.next = e.next
	e.next.prev = e.prev

	e.prev = at
	e.next = at.next
	e.prev.next = e
	e.next.prev = e

	return e
}

// Remove 删除 list 上指定的节点
func (l *List[T]) Remove(e *Element[T]) T {
	if e.list == l {
		l.remove(e)
	}
	return e.Value
}

func (l *List[T]) remove(e *Element[T]) {
	e.prev.next = e.next
	e.next.prev = e.prev
	e.next = nil // 避免内存泄漏
	e.prev = nil
	e.list = nil
	l.len--
}

// PushFront 在 list 头部插入一个节点
func (l *List[T]) PushFront(v T) *Element[T] {
	l.lazyInit()
	return l.insertValue(v, &l.root)
}

func (l *List[T]) insertValue(v T, at *Element[T]) *Element[T] {
	return l.insert(&Element[T]{Value: v}, at)
}

// insert 将 e 插入到 at 后面
func (l *List[T]) insert(e, at *Element[T]) *Element[T] {
	e.prev = at
	e.next = at.next
	e.prev.next = e
	e.next.prev = e
	e.list = l
	l.len++
	return e
}

// Back 返回 list 最后的节点
func (l *List[T]) Back() *Element[T] {
	if l.len == 0 {
		return nil
	}

	return l.root.prev
}
