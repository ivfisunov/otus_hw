package hw04_lru_cache //nolint:golint,stylecheck

type List interface {
	Len() int                          // длина списка
	Front() *listItem                  // первый элемент списка
	Back() *listItem                   // последний элемент списка
	PushFront(v interface{}) *listItem // добавить значение в начало
	PushBack(v interface{}) *listItem  // добавить значение в конец
	Remove(i *listItem)                // удалить элемент
	MoveToFront(i *listItem)           // переместить элемент в начало
}

type listItem struct {
	Value interface{}
	Next  *listItem
	Prev  *listItem
}

type list struct {
	length int
	first  *listItem
	last   *listItem
}

func (l *list) Len() int {
	return l.length
}

func (l *list) Front() *listItem {
	return l.first
}

func (l *list) Back() *listItem {
	return l.last
}

func (l *list) PushFront(v interface{}) *listItem {
	newNode := listItem{}
	newNode.Value = v
	if l.length == 0 {
		l.first, l.last = &newNode, &newNode
		l.length++
		return &newNode
	}
	newNode.Next = l.first
	l.first.Prev, l.first = &newNode, &newNode
	l.length++
	return &newNode
}

func (l *list) PushBack(v interface{}) *listItem {
	newNode := listItem{}
	newNode.Value = v
	if l.length == 0 {
		l.first, l.last = &newNode, &newNode
		l.length++
		return &newNode
	}
	newNode.Prev = l.last
	l.last.Next, l.last = &newNode, &newNode
	l.length++
	return &newNode
}

func (l *list) Remove(i *listItem) {
	if l.length == 0 {
		return
	}
	if i == l.first {
		l.first = i.Next
		l.first.Prev, i.Next = nil, nil
		l.length--
		return
	}
	if i == l.last {
		l.last = i.Prev
		l.last.Next, i.Prev = nil, nil
		l.length--
		return
	}
	i.Prev.Next, i.Next.Prev = i.Next, i.Prev
	l.length--
}

func (l *list) MoveToFront(i *listItem) {
	if l.length == 0 {
		return
	}
	value := i.Value
	l.Remove(i)
	l.PushFront(value)
}

func NewList() List {
	return &list{}
}
